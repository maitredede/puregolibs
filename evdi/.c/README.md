gcc -o evdi_example evdi_example.c -levdi -Wall
gcc main.c -o simulate_display -levdi -Wall
sudo timeout 120 ./simulate_display
