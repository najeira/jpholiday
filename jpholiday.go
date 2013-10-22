package jpholiday

import (
	"math"
	"time"
)

var JST = time.FixedZone("JST", 3600*9)

func Name(t time.Time) string {
	t = t.In(JST)
	y, m, d, w := toYMDW(t)
	name := getName(y, m, d, w)

	//振替休日
	if len(name) <= 0 {
		if 1973 <= y && w == 0 {
			yname := getYesterdayNameFromTime(t)
			if len(yname) >= 1 {
				name = "振替休日"
			}
		} else if m == 5 && d == 6 && 2007 <= y && 1 <= w && w <= 2 {
			name = "振替休日"
		}
	}
	return name
}

func toYMDW(t time.Time) (int, int, int, int) {
	return t.Year(), int(t.Month()), t.Day(), (int(t.Weekday()) + 6) % 7
}

func getYesterdayNameFromTime(t time.Time) string {
	yesterday := t.AddDate(0, 0, -1)
	y, m, d, w := toYMDW(yesterday)
	return getName(y, m, d, w)
}

func getName(y, m, d, w int) string {
	//皇室慶弔行事に伴う休日
	if y == 1959 && m == 4 && d == 10 {
		return "皇太子・明仁親王の結婚の儀"
	} else if y == 1989 && m == 2 && d == 24 {
		return "昭和天皇の大喪の礼"
	} else if y == 1990 && m == 11 && d == 12 {
		return "即位の礼正殿の儀"
	} else if y == 1993 && m == 6 && d == 9 {
		return "皇太子・徳仁親王の結婚の儀"
	}

	//国民の祝日
	if m == 1 {
		if d == 1 {
			return "元日"
		} else {
			if 1949 <= y && y <= 1999 {
				if d == 15 {
					return "成人の日"
				}
			} else if 2000 <= y {
				if 8 <= d && d <= 14 && w == 0 {
					return "成人の日"
				}
			}
		}
	} else if m == 2 {
		if 1967 <= y {
			if d == 11 {
				return "建国記念の日"
			}
		}
	} else if m == 3 {
		if 19 <= d && d <= 22 {
			if d == shunBunDay(y) {
				return "春分の日"
			}
		}
	} else if m == 4 {
		if d == 29 {
			if y <= 1988 {
				return "天皇誕生日"
			} else if y <= 2006 {
				return "みどりの日"
			} else {
				return "昭和の日"
			}
		}
	} else if m == 5 {
		if d == 3 {
			return "憲法記念日"
		} else if d == 4 {
			if 1988 <= y && y <= 2006 && 1 <= w && w <= 5 {
				return "国民の休日"
			} else if 2007 <= y {
				return "みどりの日"
			}
		} else if d == 5 {
			return "こどもの日"
		}
	} else if m == 7 {
		if 1996 <= y && y <= 2002 {
			if d == 20 {
				return "海の日"
			}
		} else if 2003 <= y {
			if 15 <= d && d <= 21 && w == 0 {
				return "海の日"
			}
		}
	} else if m == 9 {
		if 1966 <= y && y <= 2002 {
			if d == 15 {
				return "敬老の日"
			}
		} else if 2003 <= y {
			if 15 <= d && d <= 21 && w == 0 {
				return "敬老の日"
			}
		}
		if 2009 <= y && w == 1 {
			if 21 <= d && d <= 23 {
				if d+1 == shuuBunDay(y) {
					return "国民の休日"
				}
			}
		}
		if 22 <= d && d <= 24 {
			if d == shuuBunDay(y) {
				return "秋分の日"
			}
		}
	} else if m == 10 {
		if 1966 <= y && y <= 1999 {
			if d == 10 {
				return "体育の日"
			}
		} else if 2000 <= y {
			if 8 <= d && d <= 14 && w == 0 {
				return "体育の日"
			}
		}
	} else if m == 11 {
		if d == 3 {
			return "文化の日"
		} else if d == 23 {
			return "勤労感謝の日"
		}
	} else if m == 12 {
		if 1989 <= y && d == 23 {
			return "天皇誕生日"
		}
	}

	return ""
}

func shunBunDay(y int) int {
	add := 0.242194*float64(y-1980) - math.Floor(float64(y-1980)/4.0)
	val := 0.0
	if 2100 <= y && y <= 2150 {
		val = 21.8510 + add
	} else if 1980 <= y {
		val = 20.8431 + add
	} else if 1900 <= y {
		val = 20.8357 + add
	} else if 1851 <= y {
		val = 19.8277 + add
	}
	return int(math.Floor(val))
}

func shuuBunDay(y int) int {
	add := 0.242194*float64(y-1980) - math.Floor(float64(y-1980)/4.0)
	val := 0.0
	if 2100 <= y && y <= 2150 {
		val = 24.2488 + add
	} else if 1980 <= y {
		val = 23.2488 + add
	} else if 1900 <= y {
		val = 23.2588 + add
	} else if 1851 <= y {
		val = 22.2588 + add
	}
	return int(math.Floor(val))
}
