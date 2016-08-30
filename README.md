Twelete
=======

Twelete is a tool to delete historic tweets beyond the limits of the Twitter API using your downloadable tweet archive.


How To Use
----------

 1. Request your [Twitter Archive](https://twitter.com/settings/account)
 2. Download your archive, it will be a zip file
 3. Create a new [Twitter App](https://apps.twitter.com/app/new)
   - the name/description/website are not important
 4. Make a note of your Twitter App's Consumer Key, Consumer Secret, Access Token and Access Secret from "Keys and Access Tokens"
   - you may need to click "Generate My Access Token and Token Secret" to get your Access Token and Access Secret.
 5. To delete tweets a maximum of 10,000 from before January 1st 2016, run the below (run `./twelete --help` for more options):

    ./twelete --archive ~/Downloads/your-twitter-archive.zip --before 2016-01-01 --limit 10000 --pause 100ms --consumer-key your-consumer-key --consumer-secret your-consumer-secret --access-token your-access-token --access-secret your-access-secret

 6. Delete your Twitter App

