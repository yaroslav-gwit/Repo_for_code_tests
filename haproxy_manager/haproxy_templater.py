#!bin/python3
import sys
import configparser
from jinja2 import Template
import typer
import re
import os
from os.path import exists

haproxy_site_db_location = "haproxy_site_db.ini"
haproxy_config_template_location = "haproxy_config_template.cfg"
haproxy_config_location = "haproxy.cfg"

haproxy_site_db = configparser.ConfigParser()
haproxy_site_db.read(haproxy_site_db_location)

app = typer.Typer(context_settings=dict(max_content_width=800))


@app.command()
def add_site(site_address:str=typer.Argument(..., help="Your website address, for example gateway-it.com"), \
    www_redirection:bool=typer.Option(False, help="Add www to @ redirection for the domain"), \
    backend_servers:str=typer.Option(..., help="Supply a list of coma separated backend servers"), \
    http2_backend:bool=typer.Option(False, help="Use HTTP/2 on the backend"), \
    https_backend:bool=typer.Option(False, help="Use HTTPs on the backend"), \
        ):
    
    '''
    Example: ./haproxy_templater.py add-site gateway-it.com --backend-servers 1.1.1.1:443,2.2.2.2:443
    '''
    
    if site_address not in haproxy_site_db.sections():
        haproxy_site_db.add_section(site_address)
    
    with open(haproxy_site_db_location, 'w') as configfile:
        haproxy_site_db[site_address]['www_redirection'] = str(www_redirection)
        haproxy_site_db[site_address]['http2_backend'] = str(http2_backend)
        haproxy_site_db[site_address]['https_backend'] = str(https_backend)
        haproxy_site_db[site_address]['backend_servers'] = str(backend_servers)
        haproxy_site_db.write(configfile)
    
    backend_servers_list = haproxy_site_db[site_address]['backend_servers'].split(",")
    
    print("The site " + site_address + " was added to the system!")
    # service(refresh_config=True)


@app.command()
def remove_site(site_address:str=typer.Argument(..., help="The website address you would like to remove, for example gateway-it.com")):
    '''
    Example: ./haproxy_templater.py remove-site gateway-it.com
    '''
    
    if site_address in haproxy_site_db.sections():
        haproxy_site_db.remove_section(site_address)
        with open(haproxy_site_db_location, 'w') as configfile:
            haproxy_site_db.write(configfile)
    else:
        print("The site is not in our database!")
        exit(1)
    
    print("The site " + site_address + " was removed from the system!")


@app.command()
def service(refresh_config:bool=typer.Option(False, help="Regenerates the HAProxy config and reloads the service to apply the latest settings"), \
    get_certificate:str=typer.Option("None", help="Get the SSL certificate for any given domain"), \
    all_all_certificates:bool=typer.Option(False, help="Get the SSL certificates for all domains registered on this box"),
    renew_certificate:str=typer.Option("None", help="Renews the SSL certificate for any given domain"), \
    renew_all_certificates:bool=typer.Option(False, help="Renews the SSL certificates for all domains registered on this box")):
    
    '''
    Example: ./haproxy_templater.py service --refresh-config
    '''

    if refresh_config == True:
        # Load variables
        sites_list = haproxy_site_db.sections()
        # This removes [GLOBAL] from site list
        del(sites_list[0])
        
        www_sites_list = []
        for site in sites_list:
            if haproxy_site_db[site]["www_redirection"] == "True":
                www_sites_list.append(site)
        
        backend_servers = []
        for site in sites_list:
            backend_servers.append(haproxy_site_db[site]['backend_servers'])

        https_backend_list = []
        for site in sites_list:
            https_backend_list.append(haproxy_site_db[site]['https_backend'])

        http2_backend_list = []
        for site in sites_list:
            http2_backend_list.append(haproxy_site_db[site]['http2_backend'])

        if not exists("/ssl/"):
            print("Folder /ssl/ doesn't exist, please create it!")
            exit(1)
        if len(os.listdir('/ssl')) == 0:
            ssl_folder_empty = True
        else:
            ssl_folder_empty = False
        
        # Load the template
        with open(haproxy_config_template_location, 'r') as file_object:
            contents = file_object.read()
        
        # Generate the HAProxy config file
        haproxy_config_template_input = Template(contents)
        haproxy_config_template_output = haproxy_config_template_input.render(sites_list=sites_list,ssl_folder_empty=ssl_folder_empty,www_sites_list=www_sites_list, backend_servers=backend_servers, https_backend_list=https_backend_list, http2_backend_list=http2_backend_list)

        # Write HAProxy proxy config file
        with open(haproxy_config_location, 'w') as file_object:
            file_object.write(haproxy_config_template_output)

        print("The config file was recreated and the service was reloaded!")
    
    elif renew_all_certificates == True:
        print("All certificates were renewed!")
    
    elif renew_certificate != "None":
        print("The certificate was renewed!")
    
    else:
        print("ERROR! Please select at least one option: --refresh-config, --renew-certificate domain.com, --renew-all-certificates")
        print("For example: ./haproxy_templater.py service --refresh-config")
        print()
        print("To see the full list of options, execute this:")
        print("./haproxy_templater.py service --help")



if __name__ == "__main__":
    app()
