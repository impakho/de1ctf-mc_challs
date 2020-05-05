def deployClient(player_name):
    open('pipe/pipe.txt', 'a').write(player_name + '\n')