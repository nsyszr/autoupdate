--[[
Parsing algorithm of .OAM tags
- First find OAM Tags
- Check if it's a multiline comment
- Find the end of the comment. If multline run until corresponding Lua comment 
    end (double squared brackets), otherwise end of line.
- Get attribute list by split once at double colon :
- Take the right portion and find the first "
- If exists split string until quotation mark position by comma and skip the 
    last element. parse the attributes by splitting at equals sign. Left is key,
    right is value. Now find the corresponding of the quotation mark. Skip 
    escaped quotation marks. Save pos + 1 as next starting position. Take the 
    previously skipped last element, b/c it contains the left value of the
    attribute. Conact both strings and do the trick by splitting once at the 
    first equal character. Continue at next starting position by searching for
    a quotation mark and on success repeat this whole step.
]]

--[[.OAM.TEMPL: description="Demo config for setting primary IP address",
    firmware=icomos, firmware_min_version=3.1]]

--[[.OAM.PARAM: name=ipv4_address, type=string, required=true
    description="Primary IP address"]]
ipv4_address = nil
--.OAM.PARAM: name=ipv4_netmask, type=string, required=true
ipv4_netmask = nil
--.OAM.PARAM: name=ipv4_gateway, type=string, default="192.168.1.1"
ipv4_gateway = "192.168.1.1"
--.OAM.PARAM: name=test_array, type=array, default={1,2,3,4}
test_array = { 1, 2, 3, 4 }
--.OAM.PARAM: name=location, type=string, nullable=true
location = nil

cli("...ip_address=" .. ipv4_address) --asdasdad
cli("...ip_netmask=" .. ipv4_netmask)
cli("...ip_gateway=" .. ipv4_gateway)

cli("administration.profile.activate")
