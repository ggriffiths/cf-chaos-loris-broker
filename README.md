# cf-chaos-loris-broker

See [this doc](http://docs.spring.io/spring/docs/current/javadoc-api/org/springframework/scheduling/support/CronSequenceGenerator.html) for details about `schedule` option.

[Chaos Loris Repo](https://github.com/Altoros/chaos-loris).

[Chaos Loris API Docs](http://strepsirrhini-army.github.io/chaos-loris/).

### Algorithm

- **Create Service Instance:** create a scheduler, save ServiceInstance to db with scheduler and probobility from a plan
- **Delete Service Instance:** delete a scheduler and all chaoses, remove ServiceInstance from db
- **Create Service Bind:** create an app, create a chaos with , scheduler url and probo, remove ServiceInstance from db

### Note

To debug chaos-loris you can use following commands:
```
cf set-env chaos-loris JAVA_OPTS -Dlogging.level.org.springframework=TRACE
cf restart chaos-loris
```

