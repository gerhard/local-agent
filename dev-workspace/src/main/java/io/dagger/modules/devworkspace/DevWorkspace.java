package io.dagger.modules.devworkspace;

import static io.dagger.client.Dagger.dag;

import io.dagger.client.Container;
import io.dagger.client.DaggerQueryException;
import io.dagger.client.Directory;
import io.dagger.client.ReturnType;
import io.dagger.module.annotation.Function;
import io.dagger.module.annotation.Object;
import java.util.List;
import java.util.concurrent.ExecutionException;

@Object
public class DevWorkspace {
  public Container container;
  public Directory source;

  public DevWorkspace() {}

  public DevWorkspace(Directory source) {
    this.container =
        dag().container().from("alpine:3").withDirectory("/app", source).withWorkdir("/app");
  }

  /**
   * Install packages in the workspace and return exit code of the installation
   *
   * @param packages List of packages to install
   */
  @Function
  public DevWorkspace addPackages(List<String> packages)
      throws ExecutionException, DaggerQueryException, InterruptedException {
    List<String> args = new java.util.ArrayList<>();
    args.addAll(List.of("apk", "add", "--no-cache", "--update"));
    args.addAll(packages);
    return withExec(args);
  }

  /**
   * Run command and return the exit code
   *
   * @param args Command to run
   */
  @Function
  public DevWorkspace withExec(List<String> args)
      throws ExecutionException, DaggerQueryException, InterruptedException {
    this.container =
        this.container.withExec(args, new Container.WithExecArguments().withExpect(ReturnType.ANY));
    String out = this.container.stdout() + "\n" + this.container.stderr();
    int exitCode = this.container.exitCode();
    if (exitCode != 0) {
      throw new RuntimeException("Failed to execute command: " + out);
    }
    return this;
  }

  /**
   * Read the contents of a file in the workspace at the given path
   *
   * @param path Path to read the file at
   */
  @Function
  public String read(String path)
      throws ExecutionException, DaggerQueryException, InterruptedException {
    return container.file(path).contents();
  }

  /**
   * Writhe contents of a file in the workspace at the given path
   *
   * @param path Path to write the file at
   * @param contents Contents to write
   */
  @Function
  public DevWorkspace write(String path, String contents) {
    this.container = this.container.withNewFile(path, contents);
    return this;
  }

  /**
   * Remove a file in the workspace at the given path
   *
   * @param path Path to remove the file at
   */
  @Function
  public DevWorkspace Remove(String path) {
    this.container = this.container.withoutFile(path);
    return this;
  }

  /** List the files in the workspace in tree format */
  @Function
  public String tree() throws ExecutionException, DaggerQueryException, InterruptedException {
    return this.container.withExec(List.of("tree", ".")).stdout();
  }
}
