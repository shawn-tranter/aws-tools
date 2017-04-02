package main

import (
    //"fmt"
    "regexp"
    "strconv"
    "errors"
)

func parseProbability(probabilityString string) (float64, error) {
    //fmt.Println("probabilityString:", probabilityString)

    // example 35.2%
    rePercentDecimal := regexp.MustCompile("(\\d+\\.\\d+)\\%")
    matchPercentDecimal := rePercentDecimal.FindStringSubmatch(probabilityString)

    if len(matchPercentDecimal) > 0 {
        //fmt.Println("percent with decimal:", matchPercentDecimal[1])
        number, _ := strconv.ParseFloat(matchPercentDecimal[1], 64)
        return number/100, nil
    }

    // example 82%
    rePercent := regexp.MustCompile("(\\d+)\\%")
    matchPercent := rePercent.FindStringSubmatch(probabilityString)

    if len(matchPercent) > 0 {
        //fmt.Println("percent:", matchPercent[1])
        number, _ := strconv.ParseFloat(matchPercent[1], 64)
        return number/100, nil
    }

    // example 1/2 (50%)
    reFraction := regexp.MustCompile("(\\d+)\\/(\\d+)")
    matchFraction := reFraction.FindStringSubmatch(probabilityString)

    if len(matchFraction) > 0 {
        //fmt.Println("fraction:", matchFraction[1], ",", matchFraction[2])
        num, _ := strconv.ParseFloat(matchFraction[1], 64)
        den, _ := strconv.ParseFloat(matchFraction[2], 64)
        return num/den, nil
    }

    // example 0.3 (30%)
    reDecimal := regexp.MustCompile("(\\d*\\.\\d+)")
    matchDecimal := reDecimal.FindStringSubmatch(probabilityString)

    if len(matchDecimal) > 0 {
        //fmt.Println("decimal:", matchDecimal[1])
        number, _ := strconv.ParseFloat(matchDecimal[1], 64)
        return number, nil
    }

    return 0.0, errors.New("No probability could be devined")
}

func assertProbability(probability float64) (float64, error) {
    if probability <= 0 {
        return probability, errors.New("Probability less than or equal to 0")
    }

    if probability > 1 {
        return probability, errors.New("Probability greater than 1")
    }

    return probability, nil
}