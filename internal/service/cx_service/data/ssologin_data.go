package data

type AccinfoData struct {
	Msg struct {
		Fid          int `json:"fid"`
		Rosterrights int `json:"rosterrights"`
		CodeInfos    struct {
			HomeConfig struct {
				Weburl   string `json:"weburl"`
				Dwtype   int    `json:"dwtype"`
				Hometype int    `json:"hometype"`
			} `json:"homeConfig"`
		} `json:"codeInfos"`
		Boundaccount int `json:"boundaccount"`
		LoginId      int `json:"loginId"`
		CodeInfo     struct {
		} `json:"codeInfo"`
		Pic        string `json:"pic"`
		Source     string `json:"source"`
		Type       int    `json:"type"`
		Ranknum    string `json:"ranknum"`
		IsCertify  int    `json:"isCertify"`
		Uname      string `json:"uname"`
		CopyRight  int    `json:"copyRight"`
		UnitConfig struct {
		} `json:"unitConfig"`
		Schoolname     string `json:"schoolname"`
		UnitConfigInfo struct {
			HpConfig struct {
				ShowBaseHp         int `json:"showBaseHp"`
				ShowMicroServiceHp int `json:"showMicroServiceHp"`
			} `json:"hpConfig"`
			Xb int `json:"xb"`
		} `json:"unitConfigInfo"`
		UnitConfigInfos []struct {
			Fid      int `json:"fid"`
			HpConfig struct {
				ShowBaseHp         int `json:"showBaseHp"`
				ShowMicroServiceHp int `json:"showMicroServiceHp"`
			} `json:"hpConfig"`
			Xb         int    `json:"xb"`
			Schoolname string `json:"schoolname"`
		} `json:"unitConfigInfos"`
		Phone      string `json:"phone"`
		BindFanya  bool   `json:"bindFanya"`
		UpdateWay  string `json:"updateWay"`
		Name       string `json:"name"`
		Fullpinyin string `json:"fullpinyin"`
		UserConfig struct {
			Recommend struct {
				All int `json:"all"`
			} `json:"recommend"`
		} `json:"userConfig"`
		Status      int    `json:"status"`
		SwitchInfo  string `json:"switchInfo"`
		Roleid      string `json:"roleid"`
		Controlinfo struct {
			InitializedRole string `json:"initializedRole"`
			Selected        string `json:"selected"`
		} `json:"controlinfo"`
		Industry       int    `json:"industry"`
		Uid            int    `json:"uid"`
		Acttime2       string `json:"acttime2"`
		Dxfid          string `json:"dxfid"`
		Puid           int    `json:"puid"`
		Rights         int    `json:"rights"`
		NeedIntruction int    `json:"needIntruction"`
		Openid4        string `json:"openid4"`
		BindOpac       bool   `json:"bindOpac"`
		Ppfid          string `json:"ppfid"`
		AccountInfo    struct {
			CxOpac struct {
				LoginId  int    `json:"loginId"`
				Tiptitle string `json:"tiptitle"`
				LoginUrl string `json:"loginUrl"`
				BoundUrl string `json:"boundUrl"`
				Tippwd   string `json:"tippwd"`
				Tipuname string `json:"tipuname"`
			} `json:"cx_opac"`
			ImAccount struct {
				Uid       int    `json:"uid"`
				Password  string `json:"password"`
				Created   int64  `json:"created"`
				Modified  int64  `json:"modified"`
				Type      string `json:"type"`
				Uuid      string `json:"uuid"`
				Activated int    `json:"activated"`
				Username  string `json:"username"`
			} `json:"imAccount"`
			CxFanya struct {
				LoginId     int    `json:"loginId"`
				CopyRight   int    `json:"copyRight"`
				Roleid      string `json:"roleid"`
				Countrycode string `json:"countrycode"`
				Tippwd      string `json:"tippwd"`
				Result      bool   `json:"result"`
				Uid         int    `json:"uid"`
				Dxfid       string `json:"dxfid"`
				Tiptitle    string `json:"tiptitle"`
				LoginUrl    string `json:"loginUrl"`
				Schoolid    int    `json:"schoolid"`
				Time        int64  `json:"time"`
				BoundUrl    string `json:"boundUrl"`
				Tipuname    string `json:"tipuname"`
				IsCertify   int    `json:"isCertify"`
				Cxid        int    `json:"cxid"`
				CreateDate  string `json:"createDate"`
				Status      int    `json:"status"`
			} `json:"cx_fanya"`
		} `json:"accountInfo"`
		Simplepinyin         string `json:"simplepinyin"`
		Sex                  int    `json:"sex"`
		IsNewUser            int    `json:"isNewUser"`
		Studentcode          string `json:"studentcode"`
		PrivacyPolicyVersion int    `json:"privacyPolicyVersion"`
		Inputfid             int    `json:"inputfid"`
	} `json:"msg"`
	Result int `json:"result"`
}
