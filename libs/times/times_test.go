package times

import (
    "testing"
    "time"
)

type test struct {
    time    time.Time
    format  string
    strTime string
}

var testCases = []test{
    {
        time.Date(2012, 11, 22, 21, 28, 10, 0, time.Local),
        "Y-m-d H:i:s",
        "2012-11-22 21:28:10",
    },
    {
        time.Date(2012, 11, 22, 0, 0, 0, 0, time.Local),
        "Y-m-d",
        "2012-11-22",
    },
    {
        time.Date(2012, 11, 22, 21, 28, 10, 0, time.Local),
        "Y-m-d H:i:s",
        "2012-11-22 21:28:10",
    },
}

func TestFormat(t *testing.T) {
    for _, testCase := range testCases {
        strTime := Format(testCase.format, testCase.time)
        if strTime != testCase.strTime {
            t.Errorf("(expected) %v != %v (actual)", testCase.time, strTime)
        }
    }
}

func TestStrToLocalTime(t *testing.T) {
    for _, testCase := range testCases {
        time := StrToLocalTime(testCase.strTime)
        if time != testCase.time {
            t.Errorf("(expected) %v != %v (actual)", time, testCase.time)
        }
    }
}

func TestStrToTime(t *testing.T) {
    // zoneName, err := time.LoadLocation("CST")
    // if err != nil {
    //     t.Error(err)
    // }

    var testCases = []test{
        {
            time.Date(2012, 11, 22, 21, 28, 10, 0, time.Local),
            "",
            "2012-11-22 21:28:10 +0800 +0800",
        },
        {
            time.Date(2012, 11, 22, 0, 0, 0, 0, time.Local),
            "",
            "2012-11-22 +0800 +0800",
        },
        {
            time.Date(2012, 11, 22, 21, 28, 10, 0, time.FixedZone("CST", 28800)),
            "",
            "2012-11-22 21:28:10 +0800 CST",
        },
    }
    for _, testCase := range testCases {
        time := StrToTime(testCase.strTime)
        // if time != testCase.time {
        if !time.Equal(testCase.time) {
            t.Errorf("(expected) %v != %v (actual)", time, testCase.time)
        }
    }
}
