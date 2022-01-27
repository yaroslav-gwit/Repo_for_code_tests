#!/usr/bin/env python3
from jinja2 import Template
import typer
import yaml
import os
import logging
import syslog
import sys


haproxy_site_db_location = "haproxy_site_db.yml"
haproxy_config_template_location = "haproxy_config_template_yaml.cfg"
haproxy_config_location = "haproxy_yaml.cfg"


app = typer.Typer(context_settings=dict(max_content_width=800))


class ConfigOptions:
    """This class is responsible to generate, reload or test the HAProxy config"""

    def __init__(self):
        self.irony = 1
    
    def reload(self):
        return self

    def generate(self):
        yaml_db = YamlFileManipulations().read()
        template = Template(JinjaReadWrite().read())

        if os.path.exists("/ssl/"):
            if len(os.listdir("/ssl/")) != 0:
                ssl_folder_not_empty = True
        else:
            ssl_folder_not_empty = False
        
        template = template.render(ssl_folder_not_empty=ssl_folder_not_empty, yaml_db=yaml_db)
        return template

    def test(self):
        return self


class YamlFileManipulations:
    """This class is responsible for adding sites, removing sites and editing yaml site DB"""

    def __init__(self, yaml_file: str = haproxy_site_db_location, yaml_input_dict = False):
        self.yaml_file = yaml_file
        self.yaml_input_dict = yaml_input_dict

    def read(self):
        with open(self.yaml_file, 'r') as file:
            yaml_file = yaml.safe_load(file)
        return yaml_file

    def write(self):
        if self.yaml_input_dict:
            with open(self.yaml_file, 'w') as file:
                yaml.dump(self.yaml_input_dict, file)
        else:
            print("There is no input (dictionary) to work with!")
            sys.exit(119)


class SSLCerts:
    """This class is responsible for dealing with SSL certificates"""

    def __init__(self, frontend_adress = False, www_redirect = False):
        self.frontend_adress = frontend_adress
        self.www_redirect = www_redirect
    

    def new_cert_from_le(self):
        if not self.frontend_adress:
            message_ = "There was no frontend address set!"
            logging.critical(message_)
            syslog.syslog(syslog.LOG_CRIT, message_)
            sys.exit(117)

        command = "certbot certonly --standalone -d " + self.frontend_adress + " --non-interactive --agree-tos --email=slv@yari.pw --http-01-port=8888"
        subprocess.run(command, shell=True, stdout=None)

        command = "cat /etc/letsencrypt/live/" + self.frontend_adress + "/fullchain.pem /etc/letsencrypt/live/" + self.frontend_adress + "/privkey.pem > /ssl/" + self.frontend_adress + ".pem"
        subprocess.run(command, shell=True, stdout=None)

        if self.www_redirect:
            command = "certbot certonly --standalone -d www." + self.frontend_adress + " --non-interactive --agree-tos --email=slv@yari.pw --http-01-port=8888"
            subprocess.run(command, shell=True, stdout=None)

            command = "cat /etc/letsencrypt/live/www." + self.frontend_adress + "/fullchain.pem /etc/letsencrypt/live/www." + self.frontend_adress + "/privkey.pem > /ssl/www." + self.frontend_adress + ".pem"
            subprocess.run(command, shell=True, stdout=None)

        # Generate new HAProxy config
        # Create SSL check here
        # Create config check here
        # Reload the HAProxy Service here

        status = ("Success", "Failure")

        return status


    def create_self_signed(self):
        return self
    
    def retire_cert(self):
        # Copy the old cert to "archive" folder before renewal
        return self
    
    def renew_cert(self):
        # Call test_cert function to determine if renewal is needed
        # Call retire function
        # Call new_cert function
        # Return status
        return self

    def test_cert(self):
        # Check if cert exists
        # Check the date on cert
        # Return status
        return self
    
    def check_if_exist(self):
        return self


class JinjaReadWrite:
    """This class is responsible for Jinja2 template file handling"""

    def __init__(self, haproxy_config_template = haproxy_config_template_location):
        self.haproxy_config_template = haproxy_config_template
        if not os.path.exists(self.haproxy_config_template):
            message_ = "Template file doesn't exist!"
            logging.critical(message_)
            syslog.syslog(syslog.LOG_CRIT, message_)
            sys.exit(118)

    def read(self):
        with open(self.haproxy_config_template, 'r') as file:
            haproxy_config_template = file.read()
        return haproxy_config_template


@app.command()
def config(reload:bool=typer.Option(False, help="Generate, test and reload the config"), \
    generate:bool=typer.Option(False, help="Only generate new config (used for troubleshooting)"), \
    test:bool=typer.Option(False, help="Generate and test the new config (used for troubleshooting)"), \
    show:bool=typer.Option(False, help="Print out the latest config"), \
        ):
    
    '''
    Example: program 
    '''

    if (reload and generate) or (reload and test) or (generate and test):
        print("You can't use these options together!")
        sys.exit(120)

    if show:
        print(ConfigOptions().generate())


@app.command()
def site_db(add:bool=typer.Option(False, help="Generate, test and reload the config"), \
    remove:bool=typer.Option(False, help="Only generate new config (used for troubleshooting)"), \
    update:bool=typer.Option(False, help="Only generate new config (used for troubleshooting)"), \
    show:bool=typer.Option(False, help="Print out the latest config"), \
        ):

    '''
    Example: program 
    '''

    if (add and remove) or (add and update) or (remove and update):
        print("You can't use these options together!")
        sys.exit(120)

    if show:
        yaml_db = yaml.dump(YamlFileManipulations().read(), sort_keys=False)
        print(yaml_db)


@app.command()
def certificate(add:bool=typer.Option(False, help="Generate, test and reload the config"), \
    remove:bool=typer.Option(False, help="Only generate new config (used for troubleshooting)"), \
    update:bool=typer.Option(False, help="Only generate new config (used for troubleshooting)"), \
    show:bool=typer.Option(False, help="Print out the latest config"), \
        ):

    '''
    Example: program 
    '''


if __name__ == "__main__":
    app()