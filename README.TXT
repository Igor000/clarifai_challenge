Steps to build and executei the Go executable:

1 cd clarifai_challenge
2 go build
3 ./clarifai_challenge


Program logic:

1 The main() calls readHTMLPage(imageURL) which returns the list of all photos 
  listOfJpg := readHTMLPage(imageURL)

2 I use "go routine" to get Predict Information for every photo.
  This data is sent to the channel "chOut"
  To avoid overloading of the server the program sleeps for 3 sec after every 10 requests

  for i, singleJpg := range listOfJpg {
	fmt.Println("Producer loop", i, singleJpg)

	go getPredictInfo(singleJpg, predictionURL, authKey, chOut)

	if i%10 == 0 {
		fmt.Println("Producer loop sleeping i = ", i)
		time.Sleep(3000 * time.Millisecond)
	}

   }

3 I collect the results from the channel chOut and store them in the fullMap (tag => prediction data)

   for i := range listOfJpg {
	fmt.Println("Consumer loop", i)

	predictInfo := <-chOut
	singleJpg := predictInfo.singleJpg
	predictionMap := predictInfo.predictionMap

	fullMap[singleJpg] = predictionMap
   }

4 Using data in the fullMap I build a new data structure
  uniqPredictionJpgMap[tag] = PriorityQueue, which is a map with key = tag
  and value = PriorityQueue (size = 10). See file priority_queue.go for details
  Priority queue keeps top 10 tags by numerical value. It's based on Go "container/heap"

5 After we have all data the program displays the menu:

======================================
Enter <all_tags> to see the entire list of tags
Enter <tag 'tag_name'> to see top 10 for a given 'tag_name'
Enter <jpg 'image_url'> to see all tags for a given 'image_url'
Enter <quit> to quit
======================================


Enter "all_tags" to display all available tags
Enter "tag 'some_tag'" to display up to top 10 photos for this tag
Enter "jpg 'image_url' to display all tags for this image
Enter "quit" to quit


I think the data structures are pretyy efficient for Production use.
I use mostly different types of maps, list and priority queue.

I don't check for return code when I look for prediction information.
That's certainly needs to be improved

I don't check for return code when I look for prediction information.
That's certainly needs to be improved. I also need to add logic to retry
requests for prediction information. I don't have it now.
I just expect every request to be successful.
I think I hit the limit of available requests and I get now this message
"402 Payment Required"

You can change my variable on line 243
authKey := "Key 74ade0e0754147b89896566bc052db70"
to fix this problem

.
