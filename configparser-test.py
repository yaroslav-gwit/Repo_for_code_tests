import configparser
from jinja2 import Template

vm_config = configparser.ConfigParser()
host_config = configparser.ConfigParser()

vm_config.read("vm.ini")
host_config.read("host.ini")

# #_ WRITE TO CONFIG FILE _#
# host_config["GLOBAL"]["default_vm_network"] = "internal"
# with open("host.ini", "w") as configfile:
#     host_config.write(configfile)

# ip_full_range_start = "10.0.101.1"
# ip_full_range_end = "10.0.101.200"
# ip_range_start = int(ip_full_range_start.split(".")[-1])
# ip_range_end = int(ip_full_range_end.split(".")[-1])