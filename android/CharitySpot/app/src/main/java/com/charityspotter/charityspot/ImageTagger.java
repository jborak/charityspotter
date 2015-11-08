package com.charityspotter.charityspot;

import com.clarifai.api.ClarifaiClient;
import com.clarifai.api.RecognitionRequest;
import com.clarifai.api.RecognitionResult;
import com.clarifai.api.Tag;

import java.io.File;
import java.util.List;

/**
 * Created by john on 11/7/15.
 */
public class ImageTagger {

    private static final String APP_ID = "HPEW2NyFNo8vKK2zf3BRowpxs6HyCMJTcapAOEQq";
    private static final String APP_SECRET = "2WNH9BioxGSHWA_IkhjLwZyMO-Bi5po1qpLXSRPW";

    private ClarifaiClient client;

    public ImageTagger() {
        client = new ClarifaiClient(APP_ID, APP_SECRET);
    }

    public String[] getTag(File imgFile) {
        List<RecognitionResult> results = client.recognize(new RecognitionRequest(imgFile));

        int numTags = results.get(0).getTags().size();
        if (numTags == 0) {
            return null;
        }

        String[] tagStrings = new String[numTags];
        int i = 0;
        for (Tag tag : results.get(0).getTags()) {
            tagStrings[i] = tag.getName();
            i++;
        }

        return tagStrings;
    }
}
