for i in {1..10000}
do
        curl http://localhost:9090/?n=25
done
curl http://localhost:9090/finish-tracing
curl http://localhost:9090/finish

