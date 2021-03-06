package extinfo

// BasicInfoRaw contains the information sent back from the server in their raw form, i.e. no translation from ints to strings, even if possible.
type BasicInfoRaw struct {
	NumberOfClients    int    // the number of clients currently connected to the server (players and spectators)
	ProtocolVersion    int    // version number of the protocol in use by the server
	GameMode           int    // current game mode
	SecsLeft           int    // the time left until intermission in seconds
	MaxNumberOfClients int    // the maximum number of clients the server allows
	MasterMode         int    // the current master mode of the server
	Paused             bool   // wether the game is paused or not
	GameSpeed          int    // the gamespeed
	Map                string // current map
	Description        string // server description
}

// BasicInfo contains the parsed information sent back from the server, i.e. game mode and master mode are translated into human readable strings.
type BasicInfo struct {
	BasicInfoRaw
	GameMode   string // current game mode
	MasterMode string // the current master mode of the server
}

// GetBasicInfoRaw queries a Sauerbraten server at addr on port and returns the raw response or an error in case something went wrong. Raw response means that the int values sent as game mode and master mode are NOT translated into the human readable name.
func (s *Server) GetBasicInfoRaw() (BasicInfoRaw, error) {
	basicInfoRaw := BasicInfoRaw{}

	response, err := s.queryServer(buildRequest(BASIC_INFO, 0, 0))
	if err != nil {
		return basicInfoRaw, err
	}

	positionInResponse = 0

	// first int is BASIC_INFO = 1
	_ = dumpInt(response)
	basicInfoRaw.NumberOfClients = dumpInt(response)
	// next int is always 5 or 7, the number of additional attributes after the playercount and before the strings for map and description
	sevenAttributes := false
	if dumpInt(response) == 7 {
		sevenAttributes = true
	}
	basicInfoRaw.ProtocolVersion = dumpInt(response)
	basicInfoRaw.GameMode = dumpInt(response)
	basicInfoRaw.SecsLeft = dumpInt(response)
	basicInfoRaw.MaxNumberOfClients = dumpInt(response)
	basicInfoRaw.MasterMode = dumpInt(response)
	if sevenAttributes {
		if dumpInt(response) == 1 {
			basicInfoRaw.Paused = true
		}
		basicInfoRaw.GameSpeed = dumpInt(response)
	} else {
		basicInfoRaw.GameSpeed = 100
	}
	basicInfoRaw.Map = dumpString(response)
	basicInfoRaw.Description = dumpString(response)

	return basicInfoRaw, nil
}

// GetBasicInfo queries a Sauerbraten server at addr on port and returns the parsed response or an error in case something went wrong. Parsed response means that the int values sent as game mode and master mode are translated into the human readable name, e.g. '12' -> "insta ctf".
func (s *Server) GetBasicInfo() (BasicInfo, error) {
	basicInfo := BasicInfo{}

	basicInfoRaw, err := s.GetBasicInfoRaw()
	if err != nil {
		return basicInfo, err
	}

	basicInfo.BasicInfoRaw = basicInfoRaw
	basicInfo.GameMode = getGameModeName(basicInfo.BasicInfoRaw.GameMode)
	basicInfo.MasterMode = getMasterModeName(basicInfo.BasicInfoRaw.MasterMode)

	return basicInfo, nil
}
