package io.dagger.modules.reddit;

import io.dagger.client.Secret;
import io.dagger.module.annotation.Function;
import io.dagger.module.annotation.Object;
import java.util.List;

@Object
public class Reddit {

  @Object
  public static class Post {
    public String postId;
    public String title;
    public String author;
    public String body;
    public int numComments;
  }

  @Object
  public static class Comment {
    public String commentId;
    public String body;
    public String author;
    public int score;
    public String parentId;
  }

  private Secret clientId;
  private Secret clientSecret;
  private Secret username;
  private Secret password;

  public Reddit() {}

  public Reddit(Secret clientId, Secret clientSecret, Secret username, Secret password) {
    this.clientId = clientId;
    this.clientSecret = clientSecret;
    this.username = username;
    this.password = password;
  }

  @Function
  public List<Post> posts(String subreddit) throws Exception {
    RedditMonitor monitor =
        new RedditMonitor(
            clientId.plaintext(),
            clientSecret.plaintext(),
            username.plaintext(),
            password.plaintext());
    monitor.authenticate();

    return monitor.getRecentPosts(subreddit).stream()
        .map(
            p -> {
              var post = new Post();
              post.postId = p.id;
              post.title = p.title;
              post.author = p.author;
              post.body = p.body;
              post.numComments = p.numComments;
              return post;
            })
        .toList();
  }

  @Function
  public List<Comment> comments(String postId) throws Exception {
    RedditMonitor monitor =
        new RedditMonitor(
            clientId.plaintext(),
            clientSecret.plaintext(),
            username.plaintext(),
            password.plaintext());
    monitor.authenticate();

    return monitor.getCommentsForPost(postId).stream()
        .map(
            c -> {
              var comment = new Comment();
              comment.commentId = c.id;
              comment.body = c.body;
              comment.author = c.author;
              comment.score = c.score;
              comment.parentId = c.parentId;
              return comment;
            })
        .toList();
  }
}
