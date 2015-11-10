# CharitySpotter

**[CharitySpotter](www.charityspotter.com) is an open source tool devised to help charities with inventory and stored object location.**

CharitySpotter allows organizations to capture photos of their collections and automatically generate a searchable inventory.


## How it works

1. Volunteers collecting donations take pictures with their phone and push the pictures through the app.

2. CharityStopper uses the [clarifai API](http://www.clarifai.com/) which returns a list of tags based on the images of the objects.

3. Photos and tags, along with other meta-data (i.e. location of collection, storage unit, timestamp of collection, etc.) are saved on [Firebase](www.firebase.com).

4. Batch CRON jobs create an index the new objects stored, making the donations collected searchable.

5. Through the webpage (www.charityspotter.com) users and organizations can search for specific items, or browse through the different categories.

![alt text][see-snap-tag]

## What CharitySpotter brings to the table

At the **organization level** we believe CharitySpotter can help organizations keep inventories easily and make use of this inventories in order to create efficient campaigns.

Having this information available also **enhances communication between different organizations**, being able to collaborate and share information and donations as needed and establish partnerships.

At the **individual level** the tool can help connecting individuals with needs to organizations that can fill those needs.


## The Team

CharitySpotter was created as part of the [Firebase Hackathon](https://firehack.splashthat.com/) (Nov, 2015) by [John](https://github.com/jborak), [Kay](https://github.com/igweckay), [Morrison](https://github.com/codeledger) and [Will](https://github.com/WillahScott).


## Technology

#### Frontend
The [CharitySpotter](www.charitystopper.com) website was built using [bootstrap](http://getbootstrap.com/)

#### Backend
The backend is powered by [Firebase](https://www.firebase.com/) which is a platform that provides Backend-as-a-Service for real-time data that simplifies interaction between database and clients at micro-second speed. Firebase also provides the NoSQL storage of the data captured by CharitySpotter.

A [Google Compute Engine](https://cloud.google.com/compute/) runs batch processes such as the clarifai API call and the indexing on the new tagged data. GoLang scripts integrate all these features together.

#### Mobile
CharitySpotter has a working prototype Android app for object capturing. We are hoping to publish the app to Google Play soon, so stay tuned!

In the future we will make our best to release an iOS app. However, this is a side-project created by a very cool team of people that got together for the weekend, so please bear with us!


## Further Information

For further information [here](http://prezi.com/yq6hak_mjexm/?utm_campaign=share&utm_medium=copy) is the actual pitch the team presented at the event. 



[see-snap-tag]: https://raw.githubusercontent.com/jborak/charityspotter/master/img/see-snap-tag.png "Easy as 1-2-3"

