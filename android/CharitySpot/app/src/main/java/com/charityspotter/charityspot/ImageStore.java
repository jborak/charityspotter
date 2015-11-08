package com.charityspotter.charityspot;

import com.firebase.client.Firebase;

import org.json.JSONArray;
import org.json.JSONObject;

import java.io.File;
import java.util.UUID;
import java.util.Date;

/**
 * Created by john on 11/7/15.
 */
public class ImageStore {

    public class Item {
        private String uid;
        private String url;
        private String[] tags;
        private Number created;

        public Item() {
        }

        public String getUid() { return uid; }
        public String getUrl() { return url; }
        public String[] getTags() { return tags; }
        public Number getCreated() { return created; }

        public void setUid(String uid) { this.uid = uid; }
        public void setUrl(String url) { this.url = url; }
        public void setTags(String[] tags) { this.tags = tags; }
        public void setCreated(Number created) { this.created = created; }
    }

    private final static String ItemsURL = "https://charitysandbox.firebaseio.com/items.json";

    private Firebase client;

    public ImageStore() {
        client = new Firebase(ItemsURL);
    }

    public void store(File image, String[] tags) {
        String uuid = UUID.randomUUID().toString();
        Date now = new Date();

        try {
            Item item = new Item();
            item.setUid(uuid);
            item.setUrl("");
            item.setTags(tags);

            Number n = new Long(now.getTime());
            item.setCreated(n);

            client.push().setValue(item);
        } catch (Exception e) {

        }
    }
}
