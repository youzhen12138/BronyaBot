package data

type Course struct {
	Result      int    `json:"result"`
	Msg         string `json:"msg"`
	ChannelList []struct {
		Cfid     int    `json:"cfid"`
		Norder   int    `json:"norder"`
		CataName string `json:"cataName"`
		Cataid   string `json:"cataid"`
		Id       int    `json:"id"`
		Cpi      int    `json:"cpi"`
		Key      int    `json:"key"`
		Content  struct {
			Studentcount int    `json:"studentcount"`
			Chatid       string `json:"chatid"`
			IsFiled      int    `json:"isFiled"`
			Isthirdaq    int    `json:"isthirdaq"`
			Isstart      bool   `json:"isstart"`
			Isretire     int    `json:"isretire"`
			Name         string `json:"name"`
			Course       struct {
				Data []struct {
					BelongSchoolId     string `json:"belongSchoolId"`
					Coursestate        int    `json:"coursestate"`
					Teacherfactor      string `json:"teacherfactor"`
					IsCourseSquare     int    `json:"isCourseSquare"`
					CourseSquareUrl    string `json:"courseSquareUrl"`
					Imageurl           string `json:"imageurl"`
					Name               string `json:"name"`
					DefaultShowCatalog int    `json:"defaultShowCatalog"`
					Id                 int    `json:"id"`
					AppData            int    `json:"appData"`
					Schools            string `json:"schools,omitempty"`
					AppInfo            string `json:"appInfo,omitempty"`
				} `json:"data"`
			} `json:"course"`
			Roletype int    `json:"roletype"`
			Id       int    `json:"id"`
			State    int    `json:"state"`
			Cpi      int    `json:"cpi"`
			Bbsid    string `json:"bbsid"`
			IsSquare int    `json:"isSquare"`
		} `json:"content"`
		Topsign int `json:"topsign"`
	} `json:"channelList"`
	Mcode            string `json:"mcode"`
	Createcourse     int    `json:"createcourse"`
	TeacherEndCourse int    `json:"teacherEndCourse"`
	ShowEndCourse    int    `json:"showEndCourse"`
	HasMore          bool   `json:"hasMore"`
	StuEndCourse     int    `json:"stuEndCourse"`
}
