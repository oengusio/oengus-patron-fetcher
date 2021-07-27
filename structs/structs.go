package structs

// PatronDisplay // Output shown to the user
type PatronDisplay struct {
    Name string `json:"full_name"`
    // assign this somehow (or just do it in frontend?)
    //p.Url = "https://www.patreon.com/user?u=" + p.Id
    Id string `json:"id"`
}

type PatronOutput struct {
    Patrons []PatronDisplay `json:"patrons"`
}

// PatreonTokens // Credentials file
type PatreonTokens struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    TokenType string `json:"token_type"`
    Scope string `json:"scope"`
    ExpiresIn int `json:"expires_in"`
}

// PatreonMembersResponse // Patreon api responses
type PatreonMembersResponse struct {
    Data []PatreonMembersData `json:"data"`
}

type PatreonMembersData struct {
    Attributes    PatreonMembersAttribute `json:"attributes"`
    Relationships PatronRelationship      `json:"relationships"`
}

type PatronRelationship struct {
    User PatronRelationshipUser `json:"user"`
}

type PatronRelationshipUser struct {
    Data PatronRelationshipUserData `json:"data"`
}

type PatronRelationshipUserData struct {
    Id   string `json:"id"`
    Type string `json:"type"`
}

type PatreonMembersAttribute struct {
    FullName           string `json:"full_name"`
    PatronStatus       string `json:"patron_status"`
    WillPayAmountCents int    `json:"will_pay_amount_cents"`
}
