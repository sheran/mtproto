package mtproto

type UserStatus struct {
	Status    string `json:"status"`
	Online    bool   `json:"online"`
	Timestamp int32  `json:"timestamp"`
}
type UserProfilePhoto struct {
	ID         int64        `json:"id"`
	PhotoSmall FileLocation `json:"photo_small"`
	PhotoLarge FileLocation `json:"photo_large"`
}
type User struct {
	Flags                UserFlags         `json:"flags"`
	ID                   int32             `json:"id"`
	Username             string            `json:"username"`
	FirstName            string            `json:"first_name"`
	LastName             string            `json:"last_name"`
	Phone                string            `json:"phone"`
	Photo                *UserProfilePhoto `json:"photot"`
	Status               *UserStatus       `json:"status"`
	Inactive             bool              `json:"inactive"`
	Mutual               bool              `json:"mutual"`
	Verified             bool              `json:"verified"`
	Restricted           bool              `json:"restricted"`
	AccessHash           int64             `json:"access_hash"`
	BotInfoVersion       int32             `json:"bot_info_version"`
	BotInlinePlaceHolser string            `json:"bot_inline_placeholder"`
	RestrictionReason    string            `json:"restriction_reason"`
	TlUser               *TL_user          `json:"tl_user"`
}
type UserFlags struct {
	Self           bool `json:"self"`             // flags_10?true
	Contact        bool `json:"contact"`          // flags_11?true
	MutualContact  bool `json:"mutual_contact"`   // flags_12?true
	Deleted        bool `json:"deleted"`          // flags_13?true
	Bot            bool `json:"bot"`              // flags_14?true
	BotChatHistory bool `json:"bot_chat_history"` // flags_15?true
	BotNochats     bool `json:"bot_no_chats"`     // flags_16?true
	Verified       bool `json:"verified"`         // flags_17?true
	Restricted     bool `json:"restricted"`       // flags_18?true
	Min            bool `json:"min"`              // flags_20?true
	BotInlineGeo   bool `json:"bot_inline_geo"`   // flags_21?true
}

func (f *UserFlags) loadFlags(flags int32) {
	if flags&1<<10 != 0 {
		f.Self = true
	}
	if flags&1<<11 != 0 {
		f.Contact = true
	}
	if flags&1<<12 != 0 {
		f.MutualContact = true
	}
	if flags&1<<13 != 0 {
		f.Deleted = true
	}
	if flags&1<<14 != 0 {
		f.Bot = true
	}
	if flags&1<<15 != 0 {
		f.BotChatHistory = true
	}
	if flags&1<<16 != 0 {
		f.BotNochats = true
	}
	if flags&1<<17 != 0 {
		f.Verified = true
	}
	if flags&1<<18 != 0 {
		f.Restricted = true
	}
	if flags&1<<20 != 0 {
		f.Min = true
	}
	if flags&1<<21 != 0 {
		f.BotInlineGeo = true
	}
}

func (user *User) GetInputPeer() TL {
	if user.Flags.Self {
		return TL_inputPeerSelf{}
	} else {
		return TL_inputPeerUser{}
	}
}
func (user *User) GetPeer() TL {
	return TL_peerUser{
		User_id: user.ID,
	}
}
func NewUserStatus(userStatus TL) (s *UserStatus) {
	s = new(UserStatus)
	switch status := userStatus.(type) {
	case TL_userStatusEmpty:
		return nil
	case TL_userStatusOnline:
		s.Status = USER_STATUS_ONLINE
		s.Online = true
		s.Timestamp = status.Expires
	case TL_userStatusOffline:
		s.Status = USER_STATUS_OFFLINE
		s.Online = false
		s.Timestamp = status.Was_online
	case TL_userStatusRecently:
		s.Status = USER_STATUS_RECENTLY
		s.Online = false
	case TL_userStatusLastWeek:
		s.Status = USER_STATUS_LAST_WEEK
	case TL_userStatusLastMonth:
		s.Status = USER_STATUS_LAST_MONTH
	}
	return
}
func NewUserProfilePhoto(userProfilePhoto TL) (u *UserProfilePhoto) {
	u = new(UserProfilePhoto)
	switch pp := userProfilePhoto.(type) {
	case TL_userProfilePhotoEmpty:
		return nil
	case TL_userProfilePhoto:
		u.ID = pp.Photo_id
		switch big := pp.Photo_big.(type) {
		case TL_fileLocationUnavailable:
		case TL_fileLocation:
			u.PhotoLarge.DC = big.Dc_id
			u.PhotoLarge.LocalID = big.Local_id
			u.PhotoLarge.Secret = big.Secret
			u.PhotoLarge.VolumeID = big.Volume_id
		}
		switch small := pp.Photo_small.(type) {
		case TL_fileLocationUnavailable:
		case TL_fileLocation:
			u.PhotoSmall.DC = small.Dc_id
			u.PhotoSmall.LocalID = small.Local_id
			u.PhotoLarge.Secret = small.Secret
			u.PhotoSmall.VolumeID = small.Volume_id
		}
	}
	return
}
func NewUser(in TL) (user *User) {
	user = new(User)
	switch u := in.(type) {
	case TL_userEmpty:
		user.ID = u.Id
	case TL_user:
		user.TlUser = &u
		user.ID = u.Id
		user.Username = u.Username
		user.FirstName = u.First_name
		user.LastName = u.Last_name
		user.AccessHash = u.Access_hash
		user.BotInfoVersion = u.Bot_info_version
		user.BotInlinePlaceHolser = u.Bot_inline_placeholder
		user.RestrictionReason = u.Restriction_reason
		user.Phone = u.Phone
		if u.Flags&1<<5 != 0 {
			user.Photo = NewUserProfilePhoto(u.Photo)
		}
		if u.Flags&1<<6 != 0 {
			user.Status = NewUserStatus(u.Status)
		}

	default:
		//fmt.Println(reflect.TypeOf(u).String())
		return nil
	}
	return
}
