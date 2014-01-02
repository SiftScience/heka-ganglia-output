heka-ganglia-output
===================

heka-ganglia-output is a Heka plugin to send accumulated stats to Ganglia.

heka-ganglia-output is released under an MIT-style open source license -- see LICENSE for details.

Uses Jeff Buchbinder's go-gmetric library (https://github.com/jbuchbinder/go-gmetric).


Building Heka with This Plugin
------------------------------

All Heka plugins written in Go must be included in Heka at compile time (see http://hekad.readthedocs.org/en/latest/installing.html#build-include-externals).  To build Heka 0.4.2 with this plugin:

1.  Clone the Heka repositoy (at version 0.4.2):

    ```sh
    git clone https://github.com/mozilla-services/heka.git --single-branch -b v0.4.2

    # or, with older versions of git:
    git clone https://github.com/mozilla-services/heka.git --depth 1
    cd heka
    git checkout v0.4.2 -b v0.4.2
    ```

2. Add the following lines to heka/cmake/externals.cmake:

    ```sh
    git_clone(https://github.com/jbuchbinder/go-gmetric 999d61122cfc4952463759c54ddfc1f1ee32e341)
    add_external_plugin(git https://github.com/SiftScience/heka-ganglia-output master)
    ```

3.  Build Heka (see http://hekad.readthedocs.org/en/latest/installing.html#from-source for dependencies and additional details):

    ```sh
    cd heka
    source build.sh
    
    # to make a DEB package:
    make deb
    ```
