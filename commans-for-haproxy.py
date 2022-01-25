command = "certbot certonly --standalone -d " + frontend_adress + " --non-interactive --agree-tos --email=slv@yari.pw --http-01-port=8888"
subprocess.run(command, shell=True, stdout=None)

command = "cat /etc/letsencrypt/live/" + frontend_adress + "/fullchain.pem /etc/letsencrypt/live/" + frontend_adress + "/privkey.pem > /ssl/" + frontend_adress + ".pem"
subprocess.run(command, shell=True, stdout=None)

if www_redirect == "Yes":
    command = "certbot certonly --standalone -d www." + frontend_adress + " --non-interactive --agree-tos --email=slv@yari.pw --http-01-port=8888"
    subprocess.run(command, shell=True, stdout=None)

    command = "cat /etc/letsencrypt/live/www." + frontend_adress + "/fullchain.pem /etc/letsencrypt/live/www." + frontend_adress + "/privkey.pem > /ssl/www." + frontend_adress + ".pem"
    subprocess.run(command, shell=True, stdout=None)

command = "curl localhost:8002/sites/haproxy-generate-template/ > /etc/haproxy/haproxy.cfg"
subprocess.run(command, shell=True, stdout=None)

command = "systemctl reload haproxy"
subprocess.run(command, shell=True, stdout=None)