;
; BIND data file for local loopback interface
;
$TTL    604800
@       IN      SOA     dnstesting.com. admin.dnstesting.com. (
                              3         ; Serial
                         604800         ; Refresh
                          86400         ; Retry
                        2419200         ; Expire
                         604800 )       ; Negative Cache TTL
;
; name server
@                          IN      A       172.20.0.3
@                          IN      NS      ns.dnstesting.com.
ns                         IN      A       172.20.0.3

; Other A records
host                       IN      A       172.20.0.2
cname                      IN      CNAME   host.
mail                       IN      A       172.20.0.2
@                          IN      MX 10   mail
@                          IN      TXT     "txt-entry"
