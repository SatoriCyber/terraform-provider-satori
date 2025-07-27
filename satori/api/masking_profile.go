package api

import (
  "encoding/json"
  "log"
)

const MaskingProfileApiPrefix = "/api/v1/masking"

type MaskingProfile struct {
  Name        string             `json:"name"`
  Description *string            `json:"description,omitempty"`
  MaskConfigs []MaskingCondition `json:"maskConfigs"`
}

type MaskingCondition struct {
  Tag         string  `json:"tag"`
  Type        string  `json:"type"`
  Replacement *string `json:"replacement,omitempty"`
  Truncate    int     `json:"truncate,omitempty"`
  SqlFunction *string `json:"sqlFunction,omitempty"`
}

type MaskingProfileOutput struct {
  MaskingProfile
  Id string `json:"id"`
}

func (c *Client) CreateMaskingProfile(input *MaskingProfile) (*MaskingProfileOutput, error) {
  output := MaskingProfileOutput{}
  jsonInput, _ := json.Marshal(input)
  log.Printf("Going to create masking profile %s", jsonInput)
  return &output, c.postJsonForAccount(MaskingProfileApiPrefix, input, &output)
}

func (c *Client) UpdateMaskingProfile(id string, input *MaskingProfile) (*MaskingProfileOutput, error) {
  output := MaskingProfileOutput{}
  return &output, c.putJson(MaskingProfileApiPrefix, "", id, input, &output)
}

func (c *Client) GetMaskingProfile(id string) (*MaskingProfileOutput, error, int) {
  output := MaskingProfileOutput{}
  err, statusCode := c.getJsonById(MaskingProfileApiPrefix, "", id, &output)
  if statusCode == 200 {
    jsonOutput, _ := json.Marshal(output)
    log.Printf("Recieved masking profile %s", jsonOutput)
  }

  return &output, err, statusCode
}

func (c *Client) DeleteMaskingProfile(id string) error {
  return c.delete(MaskingProfileApiPrefix, id)
}
