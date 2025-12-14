[fwBasic]
status = enabled
incoming = deny
outgoing = allow
routed = disabled

[Rule0]
ufw_rule = 443/tcp LIMIT IN Anywhere
description = 
command = /usr/sbin/ufw limit in proto tcp from any to any port 443
policy = limit
direction = in
protocol = tcp
from_ip = 
from_port = 
to_ip = 
to_port = 443
iface = 
routed = 
logging = 

[Rule1]
ufw_rule = 443/tcp (v6) LIMIT IN Anywhere (v6)
description = 
command = /usr/sbin/ufw limit in proto tcp from any to any port 443
policy = limit
direction = in
protocol = tcp
from_ip = 
from_port = 
to_ip = 
to_port = 443
iface = 
routed = 
logging = 

