package com.charityspotter.charityspot;

import com.firebase.client.Firebase;

/**
 * Created by morrisonchang on 11/7/15.
 */
public class CharitySpotApplication extends android.app.Application {
    @Override
    public void onCreate() {
        super.onCreate();
        Firebase.setAndroidContext(this);
    }

}
