// Package extinfo provides easy access to the state information of a Sauerbraten game server (called 'extinfo' in the Sauerbraten source code).
package extinfo

import (
	"errors"
	"net"
	"fmt"
)

// the current position in a response ([]byte)
// needed, since values are encoded in variable amount of bytes
var positionInResponse int

// Constants describing the type of information to query for
const (
	EXTENDED_INFORMATION = 0
	BASIC_INFORMATION = 1
)

// Constants describing the type of extended information to query for
const (
	UPTIME = 0
	PLAYERSTATS = 1
	TEAM_SCORE = 2
)


// GetBasicInfo queries a Sauerbraten server at addr on port and returns the parsed response or an error in case something went wrong. Parsed response means that the int values sent as game mode and master mode are translated into the human readable name, e.g. '12' -> "insta ctf".
func GetBasicInfo(addr string, port int) (BasicInfo, error) {
	response, err := queryServer(addr, port, buildRequest(BASIC_INFORMATION, 0, 0))
	if err != nil {
		return BasicInfo{}, err
	}

	positionInResponse = 0

	basicInfo := BasicInfo{}

	// first int is BASIC_INFORMATION = 1
	_ = dumpInt(response)

	basicInfo.NumberOfClients = dumpInt(response)
	// next int is always 5, the number of additional attributes after the playercount and the strings for map and description
	//numberOfAttributes := dumpInt(response)
	_ = dumpInt(response)
	basicInfo.ProtocolVersion = dumpInt(response)
	basicInfo.GameMode = getGameModeName(dumpInt(response))
	basicInfo.SecsLeft = dumpInt(response)
	basicInfo.MaxNumberOfClients = dumpInt(response)
	basicInfo.MasterMode = getMasterModeName(dumpInt(response))
	basicInfo.Map = dumpString(response)
	basicInfo.Description = dumpString(response)

	return basicInfo, nil
}

// GetBasicInfoRaw queries a Sauerbraten server at addr on port and returns the raw response or an error in case something went wrong. Raw response means that the int values sent as game mode and master mode are NOT translated into the human readable name.
func GetBasicInfoRaw(addr string, port int) (BasicInfoRaw, error) {
	response, err := queryServer(addr, port, buildRequest(BASIC_INFORMATION, 0, 0))
	if err != nil {
		return BasicInfoRaw{}, err
	}

	positionInResponse = 0

	basicInfoRaw := BasicInfoRaw{}

	// first int is always '1'
	_ = dumpInt(response)
	basicInfoRaw.NumberOfClients = dumpInt(response)
	// next int is always 5, the number of additional attributes after the playercount and the strings for map and description
	//numberOfAttributes := dumpInt(response)
	_ = dumpInt(response)
	basicInfoRaw.ProtocolVersion = dumpInt(response)
	basicInfoRaw.GameMode = dumpInt(response)
	basicInfoRaw.SecsLeft = dumpInt(response)
	basicInfoRaw.MaxNumberOfClients = dumpInt(response)
	basicInfoRaw.MasterMode = dumpInt(response)
	basicInfoRaw.Map = dumpString(response)
	basicInfoRaw.Description = dumpString(response)

	return basicInfoRaw, nil
}

// GetUptime returns the uptime of the server in seconds.
func GetUptime(addr string, port int) (int, error) {
	response, err := queryServer(addr, port, buildRequest(EXTENDED_INFORMATION, UPTIME, 0))
	if err != nil {
		return -1, err
	}

	positionInResponse = 0

	// first int is 0
	_ = dumpInt(response)

	// next int is EXT_UPTIME = 0
	_ = dumpInt(response)

	// next int is EXT_ACK = -1
	_ = dumpInt(response)

	// next int is EXT_VERSION
	_ = dumpInt(response)

	// next int is the actual uptime
	uptime := dumpInt(response)

	return uptime, nil
}

// GetPlayerInfo returns the parsed information about the player with the given clientNum.
func GetPlayerInfo(addr string, port int, clientNum int) (PlayerInfo, error) {
	playerInfo := PlayerInfo{}
	response, err := queryServer(addr, port, buildRequest(EXTENDED_INFORMATION, PLAYERSTATS, clientNum))
	if err != nil {
		return playerInfo, err
	}

	if response[5] != 0x00 {
		// there was an error
		return playerInfo, errors.New("invalid cn")
	}

	// throw away 7 first ints (EXTENDED_INFORMATION, PLAYERSTATS, clientNum, server ACK byte, server VERSION byte, server NO_ERROR byte, server PLAYERSTATS_RESP_STATS byte)
	response = response[7:]

	positionInResponse = 0


	playerInfo.ClientNum = dumpInt(response)
	playerInfo.Ping = dumpInt(response)
	playerInfo.Name = dumpString(response)
	playerInfo.Team = dumpString(response)
	playerInfo.Frags = dumpInt(response)
	playerInfo.Flags = dumpInt(response)
	playerInfo.Deaths = dumpInt(response)
	playerInfo.Teamkills = dumpInt(response)
	playerInfo.Damage = dumpInt(response)
	playerInfo.Health = dumpInt(response)
	playerInfo.Armour = dumpInt(response)
	playerInfo.Weapon = getWeaponName(dumpInt(response))
	playerInfo.Privilege = getPrivilegeName(dumpInt(response))
	playerInfo.State = getStateName(dumpInt(response))
	// IP from next 4 bytes
	ip := response[positionInResponse:positionInResponse+4]
	playerInfo.IP = net.IPv4(ip[0], ip[1], ip[2], ip[3])

	return playerInfo, nil
}

// GetPlayerInfoRaw returns the raw information about the player with the given clientNum.
func GetPlayerInfoRaw(addr string, port int, clientNum int) (PlayerInfoRaw, error) {
	playerInfoRaw := PlayerInfoRaw{}
	response, err := queryServer(addr, port, buildRequest(EXTENDED_INFORMATION, PLAYERSTATS, clientNum))
	if err != nil {
		return playerInfoRaw, err
	}

	if response[5] != 0x00 {
		// there was an error
		return playerInfoRaw, errors.New("invalid cn")
	}

	// throw away 7 first ints (EXTENDED_INFORMATION, PLAYERSTATS, clientNum, server ACK byte, server VERSION byte, server NO_ERROR byte, server PLAYERSTATS_RESP_STATS byte)
	response = response[7:]
	
	positionInResponse = 0

	playerInfoRaw.ClientNum = dumpInt(response)
	playerInfoRaw.Ping = dumpInt(response)
	playerInfoRaw.Name = dumpString(response)
	playerInfoRaw.Team = dumpString(response)
	playerInfoRaw.Frags = dumpInt(response)
	playerInfoRaw.Flags = dumpInt(response)
	playerInfoRaw.Deaths = dumpInt(response)
	playerInfoRaw.Teamkills = dumpInt(response)
	playerInfoRaw.Damage = dumpInt(response)
	playerInfoRaw.Health = dumpInt(response)
	playerInfoRaw.Armour = dumpInt(response)
	playerInfoRaw.Weapon = dumpInt(response)
	playerInfoRaw.Privilege = dumpInt(response)
	playerInfoRaw.State = dumpInt(response)
	// IP from next 4 bytes
	ip := response[positionInResponse:positionInResponse+4]
	playerInfoRaw.IP = net.IPv4(ip[0], ip[1], ip[2], ip[3])

	return playerInfoRaw, nil
}

func GetAllPlayerInfo(addr string, port int) ([]PlayerInfo, error) {
	allPlayerInfo := []PlayerInfo{}
	response, err := queryServer(addr, port, buildRequest(EXTENDED_INFORMATION, PLAYERSTATS, -1))
	if err != nil {
		return allPlayerInfo, err
	}
	fmt.Print(response)
	return allPlayerInfo, nil
}
