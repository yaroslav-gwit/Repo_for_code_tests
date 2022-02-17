#!/bin/bash
echo "Running YUM Cleanup"
sudo yum clean all

echo "Installing EPEL"
sudo yum -y install epel-release

echo "Running YUM Update"
sudo yum -y update

echo "!!! Init script - DONE !!!"

# This exit code is here to let packer know that everything went well
exit 127