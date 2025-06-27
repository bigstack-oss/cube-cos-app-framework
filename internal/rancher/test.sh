

curl 'https://10.32.10.180:10443/v3/cloudcredentials' \
  -H 'accept: application/json' \
  -H 'accept-language: en-US,en;q=0.9,zh-TW;q=0.8,zh;q=0.7' \
  -H 'content-type: application/json' \
  -H 'cookie: R_PCS=light; R_LOCALE=en-us; R_REDIRECTED=true; CSRF=a2dc7971c4534547bfcc38a9856c550f; mellon-cookie=69e06e87a1fbf947dddcc909e3d65c80; recent_project=2a3e3f3712004a67ad3f7fbceeaeccea; sessionid=08x9gedxopls122wmcxktv951mz9ejgr; csrftoken=iv6pdfT8t8nGYWIiFm3ZGyCwv6NDIRApomYPUJCye7sZvMlPHtf2BflIFnkZnLKD; R_SESS=token-r824x:kqmqjx9fpg47l9mr5gdqsz87xwdsrqcphrqtjjs6stcgjpc8mcdwkw' \
  -H 'origin: https://10.32.10.180:10443' \
  -H 'priority: u=1, i' \
  -H 'referer: https://10.32.10.180:10443/dashboard/c/_/manager/cloudCredential/create?type=openstack' \
  -H 'sec-ch-ua: "Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "macOS"' \
  -H 'sec-fetch-dest: empty' \
  -H 'sec-fetch-mode: cors' \
  -H 'sec-fetch-site: same-origin' \
  -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36' \
  -H 'x-api-csrf: a2dc7971c4534547bfcc38a9856c550f' \
  --data-raw '{"type":"provisioning.cattle.io/cloud-credential","metadata":{"generateName":"cc-","namespace":"fleet-default"},"_name":"test","annotations":{"provisioning.cattle.io/driver":"openstack"},"openstackcredentialConfig":{"password":"test"},"_type":"provisioning.cattle.io/cloud-credential","name":"test"}' \
  --insecure




{
    "type": "provisioning.cattle.io/cloud-credential",
    "metadata": {
        "generateName": "cc-",
        "namespace": "fleet-default"
    },
    "annotations": {
        "provisioning.cattle.io/driver": "openstack"
    },
    "openstackcredentialConfig": {
        "password": "test"
    },
    "name": "test-2"
}