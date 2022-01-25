#!/usr/bin/env python3
from jinja2 import Template
import typer
from os.path import exists
import yaml
import os


haproxy_site_db_location = "haproxy_site_db.yml"
haproxy_config_template_location = "haproxy_config_template_yaml.cfg"
haproxy_config_location = "haproxy_yaml.cfg"


app = typer.Typer(context_settings=dict(max_content_width=800))


class ConfigOptions:
    """This class is responsible to generate, reload or test the HAProxy config"""

    def __init__(self):
        return self
    
    def reload(self):
        return self

    def generate(self):
        return self

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
            exit(119)


class SSLCerts:
    """This class is responsible for dealing with SSL certificates"""


class JinjaReadWrite:
    """This class is responsible for Jinja2 template file handling"""

    def __init__(self, haproxy_config_template = haproxy_config_template_location):
        self.haproxy_config_template = haproxy_config_template
    
    def read(self):
        with open(self.haproxy_config_template, 'r') as file:
            haproxy_config_template = file.read()
        return haproxy_config_template


@app.command()
def config(reload:bool=typer.Option(False, help="Generate, test and reload the config"), \
    generate:bool=typer.Option(False, help="Only generate new config (used for troubleshooting)"), \
    test:bool=typer.Option(False, help="Only generate new config (used for troubleshooting)"), \
    show:bool=typer.Option(False, help="Print out the latest config"), \
        ):
    
    '''
    Example: program 
    '''

    if (reload and generate) or (reload and test) or (generate and test):
        print("You can't use these options together!")
        exit(120)

    if generate:
        yaml_db = YamlFileManipulations().read()
        template = Template(JinjaReadWrite().read())

        yaml_db_site_list = []
        yaml_db_www_site_list = []
        
        if len(os.listdir("/ssl/")) == 0:
            ssl_folder_empty = True
        else:
            ssl_folder_empty = False

        for item in range(0, len(yaml_db["sites"])):
            yaml_db_site_list.append(yaml_db["sites"][item]["site_name"])
            
            if yaml_db["sites"][item]["www_redirection"]:
                yaml_db_www_site_list.append(yaml_db["sites"][item]["site_name"])
        
        template = template.render(yaml_db_site_list=yaml_db_site_list,
            yaml_db_www_site_list=yaml_db_www_site_list,
            ssl_folder_empty=ssl_folder_empty,
            )
        
        print(template)


@app.command()
def db(add:bool=typer.Option(False, help="Generate, test and reload the config"), \
    remove:bool=typer.Option(False, help="Only generate new config (used for troubleshooting)"), \
    update:bool=typer.Option(False, help="Only generate new config (used for troubleshooting)"), \
    show:bool=typer.Option(False, help="Print out the latest config"), \
        ):

    '''
    Example: program 
    '''

    if (add and remove) or (add and update) or (remove and update):
        print("You can't use these options together!")
        exit(120)

    if show:
        yaml_db = yaml.dump(YamlFileManipulations().read(), sort_keys=False)
        print(yaml_db)


if __name__ == "__main__":
    app()