# mila_go
mila_go is a module, one to get data from web form, which uses fiber.io assign to struct.
## Features
### Cast data
* Get data from web form   
### Validate 
* Validations for structs and individual fields based on tags.

## Usage 

* Model

```
package model

import (
	. "github.com/LangPham/mila_cast"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	UserName      string `cast:"user_name" validate:"required"`
	Password      string
	PasswordPlain string `cast:"password_plain" gorm:"-"`
	Role          string `validate:"required"`
	Credentials   []Credential
}

func (models Account) Change(c *fiber.Ctx) (exchange Exchange) {

	exchange = Cast(models, c)
	exchange.ValidateModel()

	return
}
```

* Create

```
func CreateAccount(c *fiber.Ctx) (exchange Exchange) {
	account := new(Account)
	exchange = account.Change(c)

	if exchange.Valid {
		acc := exchange.Data.(Account)
		Repo.Create(&acc)
		exchange.ResultID = strconv.Itoa(int(acc.ID))
	}
	return
}  
```

* Template
```
{{#form_for exchange "/admin/account" }}
    {{#if @f}}
        <div class="error">
            {{display @f}}
        </div>
    {{/if}}

    <div>
        <label>Username:</label>
        {{ tag_input @f "user_name" value=@d.user_name }}
    </div>

    <div>
        <p><label>Email:</label></p>
        {{ tag_input @f "email" type="email" value=@d.email }}
    </div>

    <div>
        <input type="submit" value="Save">
    </div>
{{/form_for}}

```
[Example](https://github.com/LangPham/mila_go)