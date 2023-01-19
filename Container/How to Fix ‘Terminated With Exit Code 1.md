https://komodor.com/learn/how-to-fix-container-terminated-with-exit-code-1/

## What is Exit Code 1

Exit Code 1 indicates that a container shut down, either because of an application failure or because the image pointed to an invalid file. In a Unix/Linux operating system, when an application terminates with Exit Code 1, the operating system ends the process using Signal 7, known as SIGHUP.

In Kubernetes, container exit codes can help you diagnose issues with pods. If a pod is unhealthy or frequently shuts down, you can diagnose the problem using the command `kubectl describe pod [POD_NAME]`

If you see containers terminated with Exit Code 1, you’ll need to investigate the container and its applications more closely to see what caused the failure. We’ll provide several techniques for diagnosing and debugging Exit Code 1 in containers.

## Why Do Exit Code 1 Errors Occur

Exit Code 1 means that a container terminated, typically due to an application error or an invalid reference.

An application error is a programming error in any code running within the container. For example, if a Java library is running within the container, and the library throws a compiler error, the container might terminate with Exit Code 1.

An invalid reference is a file reference in the image used to run the container, which points to a nonexistent file.

### What is Signal 7 (SIGHUP)?

In Unix or Linux operating systems, signals help manage the process lifecycle. When a container terminates with Exit Code 1, the operating system terminates the container’s process with Signal 7.

Signal 7 is also known as SIGHUP – a term that originates from POSIX-compliant terminals. In old terminals based on the RS-232 protocol, SIGHUP was a “hang up” indicating the terminal has shut down.

To send Signal 7 (SIGHUP) to a Linux process use the following command:

```
kill - HUB [processID]
```

## Diagnosing Exit Code 1Checking Exit Codes in Kubernetes

## Diagnosing Exit Code 1

When a container exits, the container engine’s command-line interface (CLI) displays a line like this. The number in the brackets is the Exit Code.

```
Exited (1)
```

To list all containers that exited with an error code or didn’t start correctly:
If you are using Docker, run `ps -la`

 **Next-Gen Kubernetes Dashboard**

Manage and troubleshoot your Kubernetes workloads across multiple clusters.



**Try now - Free Forever:**

Start free with Github

[Sign up with Google, Microsoft, or Email](https://app.komodor.com/create-free-account#mode=signUp)

To diagnose why your container exited, look at the container engine logs:

- Check if a file listed in the image specification was not found. If so, the container probably exited because of this invalid reference.
- If there are no invalid references, check the logs for a clue that might indicate which library within the container caused the error, and debug the library.

## Checking Exit Codes in Kubernetes

If you are running containers as part of a Kubernetes cluster, you can find the exit code by gathering information about a pod.

Run the `kubectl describe pod [POD_NAME]` command.

The result will look something like this:

```
Containers:
  kubedns:
    Container ID: ... 
    Image:        ...
    ...
    State:          Running
      Started:      Fri, 15 Oct 2021 12:06:01 +0800
    Last State:   Terminated
      Reason:       Error
      Exit Code:    1
      Started:      Mon, 8 Nov 2021 22:21:42 +0800
      Finished:     Mon, 8 Nov 2021 23:07:17 +0800
    Ready:          True
    Restart Count:  1
```

## DIY Troubleshooting Techniques

### 1. Delete And Recreate the Container

It is a good idea to start your troubleshooting by recreating the container. This can clean out temporary files or other transient conditions that may be causing the error. Deleting and recreating will run the container with a fresh file system.

To delete and recreate the container:

- In Docker, use the docker stop command to stop the container, then use docker rm to completely remove the container. Rebuild the container using docker run
- In Kubernetes, you can manually kill the pod that runs your container using `kubectl delete pod [pod-name]`. You can then wait for Kubernetes to automatically restart your pod (depending on your setup), or manually restart the pod using `kubectl run [pod-name] --image=[image-name]`

\2. Bashing Into a Container To Troubleshoot Applications
If your container does not use entrypoints, and you suspect Exit Code 1 is caused by an application problem, you can bash into the container and try to identify which application is causing it to exit.

To bash into the container and troubleshoot applications:

1. Bashing into the container using the following command:

   ```
   docker run -ti --rm ${image} /bin/bash
   ```

2. You should now be running in a shell within the container. Run the application you suspect is causing the problem and see if it exits

3. If the application exits, check the application’s logs and see if it exited due to an application error, and what was the error

**Note**: Another way to troubleshoot an application is to simply run the application, with the same command line, outside the container. For this to be effective, you need to have an environment similar to that inside the container on the local machine.

### 3. Experimenting With Application Parameters

Exit Code 1 is often caused by application errors that cause the application, and the entire container, to exit. If you determine that Exit Code 1 is caused by an application, you can experiment with various configuration options of the application to prevent it from exiting.

Here is a partial list of application parameters you can try:

- Allocate more memory to the application
- Run the application without special switches or flags
- Make sure that the port the application uses is exposed to the relevant network
- Change the port used by the application
- Change environment variables
- Check for compatibility issues between the application and other libraries, or the underlying operating system

### 4. Addressing the PID 1 Problem

Some Exit 1 errors are caused by the PID 1 problem. In Linux, PID 1 is the “init process” that spawns other processes and sends signals.

Ordinarily, the container runs as PID 2, immediately under the init process, and additional applications running on the containers run as PID 3, 4, etc. If the application running on the container runs as PID 2, and the container itself as PID 3, the container may not terminate correctly.

**To identify if you have a PID 1 problem**

1. Run `docker ps -a` or the corresponding command in your container engine. Identify which application was running on the failed container.

2. Rerun the failed container. While it is running, in the system shell, use the command

    

   ```
   ps -aux
   ```

    

   to see currently running processes. The result will look something like this. You can identify your process by looking at the command at the end.

   ```
   PID USER PR NI VIRT RES %CPU %MEM TIME+ S COMMAND
   1 root 20 0 1.7m 1.2m 2.0 0.5 0:05.04 S {command used to run your application}
   ```

3. Look at the PID and USER at the beginning of the failing process. If PID is 1, you have a PID 1 problem.

**Possible solutions for the PID 1 problem**

- If the container will not start, try forcing it to start using a tool like `tini` or `dumb-init`
- If you are using `docker-compose`, add the init parameter to `docker-compose.yml`
- If you are using K8s, run the container using Share Process Namespace (PID 5)

These four techniques are only some of the possible approaches to troubleshooting and solving the Exit Code 1 error. There are many possible causes of Exit Code 1 which are beyond our scope, and additional approaches to resolving the problem.

## Troubleshooting Kubernetes Exit Codes with Komodor

As a Kubernetes administrator or user, pods or containers terminating unexpectedly can be a pain and can result in severe production issues.

Exit Code 1 is a prime example of how difficult it can be to identify a specific root cause in Kubernetes because many different problems can cause the same error. The troubleshooting process in Kubernetes is complex and, without the right tools, can be stressful, ineffective, and time-consuming.

Komodor is a Kubernetes troubleshooting platform that turns hours of guesswork into actionable answers in just a few clicks. Using Komodor, you can monitor, alert and troubleshoot **`exit code 1`** event.

For each K8s resource, Komodor automatically constructs a coherent view, including the relevant deploys, config changes, dependencies, metrics, and past incidents. Komodor seamlessly integrates and utilizes data from cloud providers, source controls, CI/CD pipelines, monitoring tools, and incident response platforms.

- Discover the root cause automatically with a **timeline that tracks all changes** in your application and infrastructure.
- Quickly tackle the issue, with easy-to-follow **remediation instructions**.
- Give your entire team a way to troubleshoot **independently without escalating**.