Cloning the code on api server is not scalable Thus we make use of Docker container, where it runs in a isolated environment and then generates the output.This is because the core can be very huge Thus we make use of Amazon EC2 instance so that it is secure and highly scalable.Runs in isolated environment also achieves parallelism. The generated output (build folder) is then put in Amazon S3. 