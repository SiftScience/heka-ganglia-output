package heka_ganglia_output

import (
	"fmt"
	"github.com/jbuchbinder/go-gmetric/gmetric"
	"github.com/mozilla-services/heka/message"
	"github.com/mozilla-services/heka/pipeline"
	"net"
	"strconv"
)

// An output that expects "stat metric" type messages from a StatAccumInput, with the stat values
// emitted in the message fields (emit_in_fields = true), and sends all received stats to the
// configured Ganglia gmond server.
type GangliaOutput struct {
	gm   *gmetric.Gmetric
	conn []*net.UDPConn

	group string
}

type GangliaOutputConfig struct {
	Address string // address of the Ganglia gmond server

	Host  string // present as this hostname to Ganglia
	Spoof string // Ganglia spoof string

	Group string // Ganglia metric group
}

func (o *GangliaOutput) ConfigStruct() interface{} {
	return &GangliaOutputConfig{}
}

func (o *GangliaOutput) Init(config interface{}) (err error) {
	conf := config.(*GangliaOutputConfig)
	o.gm = &gmetric.Gmetric{Host: conf.Host, Spoof: conf.Spoof}
	udp, err := net.Dial("udp", conf.Address)
	if err != nil {
		err = fmt.Errorf("Failed to open UDP connection to '%s': %v", conf.Address, err)
		return
	}
	o.conn = []*net.UDPConn{udp.(*net.UDPConn)}
	o.group = conf.Group
	return
}

func (o *GangliaOutput) Run(runner pipeline.FilterRunner, helper pipeline.PluginHelper) (err error) {
	for pack := range runner.InChan() {
		for _, field := range pack.Message.Fields {
			// TODO(jburnim): Check for and recovery from errors in sending metric.
			o.sendMetric(field)
		}
		pack.Recycle()
	}
	return
}

func (o *GangliaOutput) sendMetric(field *message.Field) {
	// Skip the "timestamp" field.
	if *field.Name == "timestamp" {
		return
	}

	var value string
	var valueType uint32
	switch rawValue := field.GetValue().(type) {
	case int64:
		value = strconv.FormatInt(rawValue, 10)
		valueType = gmetric.VALUE_INT
	case float64:
		value = strconv.FormatFloat(rawValue, 'g', 8, 64)
		valueType = gmetric.VALUE_DOUBLE
	default:
		return
	}

	// TODO(jburnim): We currently send a metadata packet every time we send a metric value.
	// We may want to keep track of for which metrics we have sent metadata packets, and only
	// send a repeat metadata packet after, e.g., 5 minutes have passed.

	// TODO(jburnim): Should we provide an option to filter out zero-valued metrics here?

	o.gm.SendMetricPackets(*field.Name, value, valueType, "", gmetric.SLOPE_UNSPECIFIED, 60, 0,
		o.group, gmetric.PACKET_BOTH, o.conn)
	return
}

func init() {
	pipeline.RegisterPlugin("GangliaOutput", func() interface{} {
		return new(GangliaOutput)
	})
}
