#!/bin/bash

# ------------------------------
# prints a message as soon as 
# ZITADEL is ready
# ------------------------------

be_status=""
fe_status=""

while [[ $be_status -ne 200 || $fe_status -ne 200 ]]; do
    sleep 5
    be_status=$(curl -s -o /dev/null -I -w "%{http_code}" localhost:${BE_PORT}/clientID)
    fe_status=$(curl -s -o /dev/null -I -w "%{http_code}" localhost:${FE_PORT}/assets/environment.json)
    echo "backend (${be_status}) or frontend (${fe_status}) not ready yet"
done

echo -e "++=======================================================================================++
||                                                                                       ||
|| ZZZZZZZZZZZZ II TTTTTTTTTTTT       AAAA       DDDDDD     EEEEEEEEEE LL                ||
||          ZZ  II      TT           AA  AA      DD    DD   EE         LL                ||
||        ZZ    II      TT          AA    AA     DD      DD EE         LL                ||
||      ZZ      II      TT         AA      AA    DD      DD EEEEEEEE   LL                ||
||    ZZ        II      TT        AAAAAAAAAAAA   DD      DD EE         LL                ||
||  ZZ          II      TT       AA          AA  DD    DD   EE         LL                ||
|| ZZZZZZZZZZZZ II      TT      AA            AA DDDDDD     EEEEEEEEEE LLLLLLLLLL        ||
||                                                                                       ||
||                                                                                       ||
|| SSSSSSSSSS TTTTTTTTTTTT       AAAA       RRRRRRRR TTTTTTTTTTTT  EEEEEEEEEE DDDDDD     ||
|| SS              TT           AA  AA      RR    RR      TT       EE         DD    DD   ||
||  SS             TT          AA    AA     RR    RR      TT       EE         DD      DD ||
||   SSSSSS        TT         AA      AA    RRRRRRRR      TT       EEEEEEEE   DD      DD ||
||        SS       TT        AAAAAAAAAAAA   RRRR          TT       EE         DD      DD ||
||         SS      TT       AA          AA  RR  RR        TT       EE         DD    DD   ||
|| SSSSSSSSSS      TT      AA            AA RR    RR      TT       EEEEEEEEEE DDDDDD     ||
||                                                                                       ||
++=======================================================================================++"
