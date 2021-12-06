#include <stdio.h>
#include <unistd.h>
#include <mpi.h>
#include <sys/time.h>
int main(int argc, char** argv){
    int node;
    int i, t;
    MPI_Init_thread(&argc, &argv, MPI_THREAD_MULTIPLE, &t);
    MPI_Comm_rank(MPI_COMM_WORLD, &node);
    char *data = "hello this is a message";
    int l = 23;

    MPI_Message msg;
    MPI_Status status;

    int n = atoi(argv[1]);

    int count;
    double t0 = MPI_Wtime();
    for (i = 0; i < n; i++) {
        if (node == 0) {

	    MPI_Send(data, l, MPI_BYTE, 1, 0, MPI_COMM_WORLD);
	    MPI_Mprobe(1, 0, MPI_COMM_WORLD, &msg, &status);
	    MPI_Get_count(&status, MPI_BYTE, &count);
	    char* buffer = malloc(count);

            MPI_Mrecv(buffer, count, MPI_BYTE, &msg, MPI_STATUS_IGNORE);
        } else {
	    MPI_Mprobe(0, 0, MPI_COMM_WORLD, &msg, &status);
	    MPI_Get_count(&status, MPI_BYTE, &count);
	    char* buffer = malloc(count);

	    MPI_Mrecv(buffer, l, MPI_BYTE, &msg, &status);
	    MPI_Send(data, l, MPI_BYTE, 0, 0, MPI_COMM_WORLD);
        }
    }
    double t1 = MPI_Wtime();
    printf("elapsed = %f s, average = %f us\n", t1 - t0, (t1-t0)/n*1e6);
    MPI_Barrier(MPI_COMM_WORLD);
    MPI_Finalize();
}
