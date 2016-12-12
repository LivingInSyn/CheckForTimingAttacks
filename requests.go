package main

import (
    "fmt"
    "time"
    "net/http"
    "math"
)

func main() {
    var numQueries int = 5
    var url string = "http://localhost:8080"
    
    //build a http client
    client := &http.Client {}

    //build the response arrays
    //var goodTimes [numQueries]int64
    var goodTimes = make([]int64,numQueries,numQueries)
    //var badTimes [numQueries]int64
    var badTimes = make([]int64,numQueries,numQueries)
    
    //build the cookies
    badCookie := http.Cookie{Name: "logged_in_user", Value: "admin--8wXqJC4wakVy1ixl+bWsdbuiDR4="}
    goodCookie := http.Cookie{Name: "logged_in_user", Value: "admin--zzzzzzzzzzzzzzzzzzzzzzzzzzz="}

    //run the good cookie tests
    for i:=0; i<numQueries; i++ {
        goodTimes[i] = do_request(client,&goodCookie,url)
        //fmt.Printf("The call took %v to run.\n", goodTimes[i])
    }
    //run the bad cookie tests
    for i:=0; i<numQueries; i++ {
        badTimes[i] = do_request(client,&badCookie,url)
        //fmt.Printf("The call took %v to run.\n", badTimes[i])
    }


    goodTimes = filter_by_std_dev(goodTimes[:])
    badTimes = filter_by_std_dev(badTimes[:])

    write_results(goodTimes[:],badTimes[:])

    //fmt.Printf("The std dev of good nums is: %v", calc_std_dev(goodTimes[:]))
}

func do_request(client *http.Client, cookie *http.Cookie, url string) int64{
    req, _ := http.NewRequest("GET", url, nil)
    req.AddCookie(cookie)
    t0 := time.Now()
    /*resp, err := client.Do(req)
    if(err != nil) {
        log.Fatal(err)
    } else {
        fmt.Printf(resp.Status+"\n")
    }*/    
    client.Do(req)
    t1 := time.Now()
    return int64(t1.Sub(t0))
}

//this calculates the mean and std deviation for the population
func calc_std_dev(array []int64) (float64,float64) {
    //calculate the mean
    var total int64 = 0
    for i:=0;i<len(array);i++{
        total = total + array[i]        
    }
    var mean float64 = float64(total/int64(len(array)))
    
    //now sum all of the points minus the mean
    var subtotal float64 = 0
    for i:=0;i<len(array);i++{
        subtotal = subtotal + math.Pow((float64(array[i]) - mean),2)
    }
    //multiple by 1/n
    var std_dev float64 = (float64(1)/float64(len(array))) * subtotal
    std_dev = math.Sqrt(std_dev)
    return mean,std_dev
}
//remove everything more than 2 std_dev away from the mean
func filter_by_std_dev(array []int64) []int64 {
    var filtered []int64

    mean, std_dev := calc_std_dev(array)
    max_val := (2*std_dev) + mean
    min_val := (2*std_dev) - mean

    for i:=0;i<len(array);i++{
        if( float64(array[i]) < max_val && float64(array[i]) > min_val ) {
            filtered = append(filtered,array[i])
        }
    }
    return filtered
}
//write the results to a csv
func write_results(goodArray []int64, badArray []int64){
    goodLen := len(goodArray)
    badLen := len(badArray)
    //figure out the array indexes
    upTo := 0
    stayBelow := 0
    longest := "good"
    if(badLen < goodLen){
        upTo = goodLen
        stayBelow = badLen
        longest = "good"
    } else {
        upTo = badLen
        stayBelow = goodLen
        longest = "bad"
    }    
    //there MUST be a cleaner way to write this
    fmt.Printf("Trial,Good Cookie,Bad Cookie\n")
    for i:=0;i<upTo;i++{
        if(i < stayBelow){
            fmt.Printf("%v,%v,%v\r\n",i,goodArray[i],badArray[i])
        } else{
            if(longest == "good") {
                fmt.Printf("%v,%v,\r\n",i,goodArray[i])
            } else {
                fmt.Printf("%v,,%v\r\n",i,badArray[i])
            }
        }
    }

}

