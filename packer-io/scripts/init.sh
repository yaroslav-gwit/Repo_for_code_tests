#!/bin/bash
echo "Running YUM Cleanup"
sudo yum clean all

echo "Installing EPEL"
sudo yum -y install epel-release

echo "Running YUM Update"
sudo yum -y update

echo "Installing ICR related check daemons"
sudo yum -y install wget

sudo wget https://github.com/yaroslav-gwit/system-checks/releases/download/0.01-alpha/reboot_after_kern_update -O /usr/bin/reboot_after_kern_update
sudo chmod +x /usr/bin/reboot_after_kern_update

sudo wget https://github.com/yaroslav-gwit/system-checks/releases/download/0.01-alpha/update_checker -O /usr/bin/update_checker
sudo chmod +x /usr/bin/update_checker
/usr/bin/update_checker > /var/log/update_checker.log

cat << EOF | cat >> /etc/profile
echo ""
echo "---> Update status <---"
cat /var/log/update_checker.log
echo ""
echo "---> Kernel status <---"
/usr/bin/reboot_after_kern_update
EOF

cat << EOF | cat >> /etc/crontab
@reboot root sleep 30 && /usr/bin/update_checker > /var/log/update_checker.log
0 */4 * * * root /usr/bin/update_checker > /var/log/update_checker.log
EOF

echo "!!! Init script - DONE !!!"

# This exit code is here to let packer know that everything went well
exit 127