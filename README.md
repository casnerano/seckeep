# SecKeep — менеджер секретных данных

[![codecov](https://codecov.io/gh/casnerano/seckeep/branch/master/graph/badge.svg?token=0FNYRKK145)](https://codecov.io/gh/casnerano/seckeep)
[![unit tests](https://github.com/casnerano/seckeep/actions/workflows/unit-tests.yml/badge.svg)](https://github.com/casnerano/seckeep/actions/workflows/unit-tests.yml)
[![lint tests](https://github.com/casnerano/seckeep/actions/workflows/lint-tests.yml/badge.svg)](https://github.com/casnerano/seckeep/actions/workflows/lint-tests.yml)

## Usage

### Client

```bash
./seckeep account sign-up --login="ivan" --password="1234" -n "Ivanov Ivan"
./seckeep account sign-in --login="ivan" --password="1234"

./seckeep data create credential --login="javascript" --password="null-is-object???" --meta="For e-mail account" --meta="Work account"
./seckeep data create card --number="4012888888881881" --month-year="06.28" --owner="Ivan Ivanov" --cvv="732" --meta="My debit visa card"
./seckeep data create text --value="My secret plan to develop a new JS-library." --meta="Secret plan"
./seckeep data create text --value="Not a bug, but a feature." --meta="My list of aphorisms"
./seckeep data create document --file="./Makefile" --meta="Example doc"
./seckeep data list

./seckeep data update --index N
./seckeep data delete --index N
./seckeep data read   --index N
```

### Dev-run
```bash
make project-init
make project-run
make load-example
```
