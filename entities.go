package main

/// <editor-fold name="UserOutput">
type PatronDisplay struct {
    Name string `json:"full_name"`
    // assign this somehow (or just do it in frontend?)
    //p.Url = "https://www.patreon.com/user?u=" + p.Id
    Id string `json:"id"`
}

type PatronOutput struct {
    Patrons []PatronDisplay `json:"patrons"`
}

/// </editor-fold>

/// <editor-fold name="Refresh token response">
type PatreonTokens struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

/// </editor-fold>

/// <editor-fold name="Patrons response">
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

/// </editor-fold>
