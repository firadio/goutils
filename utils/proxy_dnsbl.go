package utils

import (
	"strings"
)

/*
type DnsBlInfo struct {
	lStatusDomainList map[string]int
}
*/

func GetBlackByIpAddrWithProxy(ipaddr string, ProxyIpport string) map[string][]string {
	dnsbl := []string{
		"bl.score.senderscore.com",
		"bl.mailspike.net",
		"bl.spameatingmonkey.net",
		"b.barracudacentral.org",
		"bl.deadbeef.com",
		"bl.emailbasura.org",
		"bl.spamcannibal.org",
		"bl.spamcop.net",
		"blackholes.five-ten-sg.com",
		"blacklist.woody.ch",
		"bogons.cymru.com",
		"cbl.abuseat.org",
		"cdl.anti-spam.org.cn",
		"combined.abuse.ch",
		"combined.rbl.msrbl.net",
		"db.wpbl.info",
		"dnsbl-1.uceprotect.net",
		"dnsbl-2.uceprotect.net",
		"dnsbl-3.uceprotect.net",
		"dnsbl.inps.de",
		"dnsbl.sorbs.net",
		"drone.abuse.ch",
		//"drone.abuse.ch",
		"duinv.aupads.org",
		"dul.dnsbl.sorbs.net",
		"dul.ru",
		"dyna.spamrats.com",
		"dynip.rothen.com",
		"http.dnsbl.sorbs.net",
		"images.rbl.msrbl.net",
		"ips.backscatterer.org",
		"ix.dnsbl.manitu.net",
		"korea.services.net",
		"misc.dnsbl.sorbs.net",
		"noptr.spamrats.com",
		"ohps.dnsbl.net.au",
		"omrs.dnsbl.net.au",
		"orvedb.aupads.org",
		"osps.dnsbl.net.au",
		"osrs.dnsbl.net.au",
		"owfs.dnsbl.net.au",
		"owps.dnsbl.net.au",
		"pbl.spamhaus.org",
		"phishing.rbl.msrbl.net",
		"probes.dnsbl.net.au",
		"proxy.bl.gweep.ca",
		"proxy.block.transip.nl",
		"psbl.surriel.com",
		"rbl.interserver.net",
		"rdts.dnsbl.net.au",
		"relays.bl.gweep.ca",
		"relays.bl.kundenserver.de",
		"relays.nether.net",
		"residential.block.transip.nl",
		"ricn.dnsbl.net.au",
		"rmst.dnsbl.net.au",
		"sbl.spamhaus.org",
		"short.rbl.jp",
		"smtp.dnsbl.sorbs.net",
		"socks.dnsbl.sorbs.net",
		"spam.abuse.ch",
		"spam.dnsbl.sorbs.net",
		"spam.rbl.msrbl.net",
		"spam.spamrats.com",
		"spamrbl.imp.ch",
		"t3direct.dnsbl.net.au",
		"tor.dnsbl.sectoor.de",
		"torserver.tor.dnsbl.sectoor.de",
		"ubl.lashback.com",
		"ubl.unsubscore.com",
		"virbl.bit.nl",
		"virus.rbl.jp",
		"virus.rbl.msrbl.net",
		"web.dnsbl.sorbs.net",
		"wormrbl.imp.ch",
		"xbl.spamhaus.org",
		"zen.spamhaus.org",
		"zombie.dnsbl.sorbs.net",
	}
	//lStatusList := make([]int, 4)
	flg := make(chan []string)
	for _, sDomain := range dnsbl {
		go GetBlackByDomain(flg, sDomain, ipaddr, ProxyIpport)
	}
	mResult := map[string][]string{}
	for _, _ = range dnsbl {
		lInfo := <-flg
		sDomain := lInfo[0]
		sStatus := lInfo[1]
		mResult[sStatus] = append(mResult[sStatus], sDomain)
	}
	return mResult
}

func GetBlackByDomain(flg chan []string, sDomain string, ipaddr string, ProxyIpport string) {
	//defer close(flg)
	sUrl := "http://f.vision/index.php/black/" + sDomain + "/" + ipaddr
	bResult, err := ProxySocks5(ProxyIpport, sUrl)
	if err != nil {
		//fmt.Println(err)
		flg <- []string{sDomain, "ERROR"}
		return
	}
	s := string(bResult)
	//fmt.Println(s)
	if strings.Contains(s, "\"OK\"") {
		//lStatusList[1]++
		flg <- []string{sDomain, "OK"}
	} else if strings.Contains(s, "\"LISTED\"") {
		//lStatusList[2]++
		flg <- []string{sDomain, "LISTED"}
	} else {
		//lStatusList[3]++
		flg <- []string{sDomain, "UNKNOWN"}
	}
}
