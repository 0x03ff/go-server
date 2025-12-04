subjects = ["clc", "elc", "comp", "eie", "eng"]
password_suffixes = ["Password", "password"]
recovery_suffixes = ["Recover", "recover"]


# username: must tempuser
# password : ^(?:clc|elc|comp|eie|eng)[0-9]{4}(?:Password|password)$
# recover  : ^(?:clc|elc|comp|eie|eng)[0-9]{4}(?:Recover|recover)$

with open("./files/fuzz_payloads.txt", "w") as f:
    for subject in subjects:
        for num in range(0, 10000):
            num_str = f"{num:04d}"
            for p_suffix in password_suffixes:
                password = f"{subject}{num_str}{p_suffix}"
                if 8 <= len(password) <= 20:
                    for r_suffix in recovery_suffixes:
                        recover = f"{subject}{num_str}{r_suffix}"
                        if 6 <= len(recover) <= 20:
                            f.write(f'{{"username":"tempuser","password":"{password}","recover":"{recover}"}}\n')
