#include <stdio.h>
#include <unistd.h>
#include <mpi.h>
#include <sys/time.h>
int main(int argc, char** argv){
    int node;
    int t;
    MPI_Init_thread(&argc, &argv, MPI_THREAD_MULTIPLE, &t);
    MPI_Comm_rank(MPI_COMM_WORLD, &node);
    char *data = "hello there";
    int l = 11;
    struct timeval  tv1, tv2;

    MPI_Message msg;
    MPI_Status status;


    if (node == 0) {
	    printf("node 0 warming up\n");
	    fflush(stdout);
	    printf("node 0 warmup sending\n");
	    MPI_Send(data, l, MPI_BYTE, 1, 0, MPI_COMM_WORLD);
	    printf("node 0 warmup mprobing\n");
	    MPI_Mprobe(1, 0, MPI_COMM_WORLD, &msg, &status);
	    int count;
	    MPI_Get_count(&status, MPI_BYTE, &count);
	    char* buffer = malloc(count);

	    printf("node 0 warmup recving\n");
            MPI_Mrecv(buffer, count, MPI_BYTE, &msg, MPI_STATUS_IGNORE);
	    printf("node 0 sending\n");
	    fflush(stdout);
	    gettimeofday(&tv1, NULL);
	    MPI_Send(data, l, MPI_BYTE, 1, 0, MPI_COMM_WORLD);
	    MPI_Mprobe(1, 0, MPI_COMM_WORLD, &msg, &status);
	    MPI_Get_count(&status, MPI_BYTE, &count);
	    char *buffer2 = malloc(count);

            MPI_Mrecv(buffer2, count, MPI_BYTE, &msg, MPI_STATUS_IGNORE);
	    gettimeofday(&tv2, NULL);
	    double cpu_time_used;
	    printf ("Total time = %f us\n",
         (double) (tv2.tv_usec - tv1.tv_usec) +
         (double) (tv2.tv_sec - tv1.tv_sec)   * 1000000);
	    fflush(stdout);
    } else {
	    printf("node 1 warmup mprobing\n");
	    fflush(stdout);
	    MPI_Mprobe(0, 0, MPI_COMM_WORLD, &msg, &status);
	    printf("node 1 warmup recving\n");
	    int count;
	    MPI_Get_count(&status, MPI_BYTE, &count);
	    char* buffer = malloc(count);

	    MPI_Mrecv(buffer, l, MPI_BYTE, &msg, &status);
	    printf("node 1 warmup sending\n");
	    MPI_Send(data, l, MPI_BYTE, 0, 0, MPI_COMM_WORLD);

	    
	    MPI_Mprobe(0, 0, MPI_COMM_WORLD, &msg, &status);
	    MPI_Get_count(&status, MPI_BYTE, &count);
	    char* buffer2 = malloc(count);

	    MPI_Mrecv(buffer2, l, MPI_BYTE, &msg, &status);
	    MPI_Send(data, l, MPI_BYTE, 0, 0, MPI_COMM_WORLD);
    }
    MPI_Barrier(MPI_COMM_WORLD);
    MPI_Finalize();
}
