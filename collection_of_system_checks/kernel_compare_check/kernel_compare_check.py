#!/usr/bin/env python3
from natsort import natsorted
import subprocess

command = "ls -1 /boot/ | grep vmlinuz"
shell_command_ls = subprocess.check_output(command, shell=True)
shell_command_ls = shell_command_ls.decode("utf-8").split()

command = "uname -r"
shell_command_uname = subprocess.check_output(command, shell=True)
shell_command_uname = shell_command_uname.decode("utf-8").split()

latest_installed_kernel = natsorted(shell_command_ls)[-1][8:]
running_kernel = shell_command_uname[-1]

if latest_installed_kernel != running_kernel:
    print("Your OS needs a reboot to apply the latest kernel update! You are running: " + running_kernel + ". The latest installed is: " + latest_installed_kernel + ".")
else:
    print("All good! Your OS is running the latest kernel: " + latest_installed_kernel + ".")