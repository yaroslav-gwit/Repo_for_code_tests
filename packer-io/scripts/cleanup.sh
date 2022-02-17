#!/bin/bash
echo "Running install htop"
yum -y install htop

echo "!!! Cleanup script DONE !!!"

# This exit code is here to let packer know that everything went well
exit 127