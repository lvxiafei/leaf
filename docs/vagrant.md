# Running on Vagrant

If you have [Vagrant](https://www.vagrantup.com/) installed, you can run the
above example with the following commands.

1. In a terminal (terminal 1), bring up the Vagrant box:
   ```console
   $ vagrant up
   ```
   This will take a few minutes to download and provision the box.

2. Connect to the Vagrant box:
   ```console
   $ vagrant ssh
   ```

3. Build `leaf`:
   ```console
   $ cd /leaf
   $ make
   ```

4. Run `leaf`:
   ```console
   $ sudo ./leaf -H 1.1.1.1
   ```

5. In a new terminal (terminal 2), connect to the Vagrant box:
   ```console
   $ vagrant ssh
   ```

6. In terminal 2, run `curl` to generate some traffic to 1.1.1.1:
   ```console
   $ curl 1.1.1.1
   ```
   Observe the output of `leaf` in terminal 1.

7. In terminal 2, add an `iptables` rule to block traffic to 1.1.1.1:
   ```console
   $ sudo iptables -t filter -I OUTPUT 1 -m tcp --proto tcp --dst 1.1.1.1/32 -j DROP
   ```

8. In terminal 2, run `curl` to generate some traffic to 1.1.1.1:
   ```console
   $ curl 1.1.1.1
   ```
   Observe the output of `leaf` in terminal 1.

9. To clean up, press `Ctrl+C` to terminate `leaf` in terminal 1, exit both
   shells, and run:
   ```console
   $ vagrant destroy
   ```

