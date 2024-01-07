# aws_integration_with_golang

The main idea of this repository is to improve my AWS knowledge and experience using Golang as the main language of the service(s).

Before we start setting up AWS and building the zip file of the project, it is necessary to determine the `dynamDB` tablename in `main.go` file. In my case, it is `aws-dynamodb-users`.

To create a zip file for our project, you need to go to cmd: `cd cmd/` and run the command: `GOARCH=amd64 GOOS=linux go build main.go`. After that, move the `binary main file` to the `build` directory: `mv main ../build`. It is necessary to go back to the root of the project and use the command `zip -jrm build/main.zip build/main` to create a zip file.

The first action you need to do is authorize yourself to the `AWS management console` (https://eu-north-1.console.aws.amazon.com/console/home?region=eu-north-1#) and create a lambda function.

I used the name of the lambda function as `go-serverless-project1`, but of course you can use whatever name you want.

In `runtime`, it is necessary to choose Go 1.x, and I set the `architecture` to x86_64. In section `Change default execution role` I set `Create a new role from AWS policy templates` and `role name` I wrote `go-serverless-project1-executor`. In `Policy templates` I prefer `Simple microservice permissions`.

After creating the lambda function, it is necessary to change the handler in `Runtime settings`. I changed from the usual `hello` to `main`. And the last step in setting our lambda function is to upload a zip file from the `build` directory in the `code sourse` section. Well done, the lambda function is ready.

Let's start setting `dynamoDB`. Also choose `dynamoDB` in the `AWS management console` and press the button `create table`. In the table name you need to write the name that you have in your `.env` file. In my case it is `aws-dynamodb-users`. It is very significant, that name in your project absolutely matches with the name on AWS service. In the `partion key` let's choose `email` as `email` is primary key for this table. In the `table setting` I chose `Default settings`, but of course you can always customize settings as well as you want.

Now, the last step that makes everything work is to create the `API Gateway`. Go to `API Gateway` and choose an API type. I choose the classic `REST API` and press the button `Build`. In the `API details`, select `New API` and give a name to the API. I used `go-serverless-roject1` and pressed the button `Create API`. We need to create the method for the API. Choose `Create method`, in the `Method type`, select `ANY`, and in the `Integration type`, select `Lambda function`. Don't forget to turn `Lambda proxy integration` on and choose the lambda function with the name that we created earlier. Well done! The method has been created.

Now we can deploy our `API`. In `Stage` I selected `*New stage*`, and in `stage name` I used `staging`. After deploying, we get `Invoke URL` and can fully test it out!

There are our `API` commands, and the first method that we check out is `CREATE`. At the end of this curl command, you need to use your `URL` instead of just `https://.../staging`.

`POST create a user request`
```
curl --header "Content-Type: appliction/json" --request POST --data '{"email": "german@gmail.com", "first_name":"German", "last_name":"German"}' https://.../staging
```

`POST create a user respose`
```
{"id":"e926f1b0-f4ba-4be8-a5a0-a02c6f879ceb","first_name":"German","last_name":"German","email":"german@gmail.com"}
```

`GET all users request`
```
curl -X GET https://.../staging
```

`GET all users respose`
```
{"id":"e926f1b0-f4ba-4be8-a5a0-a02c6f879ceb","first_name":"German","last_name":"German","email":"german@gmail.com"}
```

`GET the user request`
```
curl -X GET https://.../staging\?email\=german@gmail.com
```

`GET the user respose`
```
{"id":"e926f1b0-f4ba-4be8-a5a0-a02c6f879ceb","first_name":"German","last_name":"German","email":"german@gmail.com"}
```
 
`PUT update the user request`
```
curl --header "Content-Type: appliction/json" --request PUT --data '{"email": "german@gmail.com", "first_name":"John", "last_name":"John"}' https://.../staging
```

`PUT update the user respose`
```
{"id":"e926f1b0-f4ba-4be8-a5a0-a02c6f879ceb","first_name":"German","last_name":"German","email":"german@gmail.com"}
```

`DELETE the user request`
```
curl -X DELETE https://.../staging\?email\=german@gmail.com
```

`DELETE the user respose`
```
null
```


