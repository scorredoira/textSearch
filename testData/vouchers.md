# Vouchers Models

## Product filter

**Table:** `voucherProduct`

### Endpoints\n\n#### Retrieve by ID\n```
GET https://[host]/api/model/voucherProduct/[id]
```

#### Retrieve many\n```
GET https://[host]/api/model/voucherProduct
```

#### Create\n```
POST https://[host]/api/model/voucherProduct
```

#### Update\n```
POST https://[host]/api/model/voucherProduct/[id]
```

#### Delete\n```
DELETE https://[host]/api/model/voucherProduct/[id]
```

### Fields\n\n| Name | Label | Type | Required | Readonly |\n|------|-------|------|----------|----------|\n| `id` | ID | int | no | yes |
| `type` | Type | int | yes | no |
| `product` | Product | int | yes | no |
| `ratio` | Ratio | float | no | no |
| `created` | Creation date | datetime | no | yes |
| `createUser` | Created by | int | no | yes |
| `updated` | Update date | datetime | no | yes |
| `updateUser` | Updated by | int | no | yes |
| `deleted` | Removed | bool | no | yes |


## Transaction

**Table:** `voucherTransaction`

### Endpoints\n\n#### Retrieve by ID\n```
GET https://[host]/api/model/voucherTransaction/[id]
```

#### Retrieve many\n```
GET https://[host]/api/model/voucherTransaction
```

#### Create\n```
POST https://[host]/api/model/voucherTransaction
```

#### Update\n```
POST https://[host]/api/model/voucherTransaction/[id]
```

#### Delete\n```
DELETE https://[host]/api/model/voucherTransaction/[id]
```

### Fields\n\n| Name | Label | Type | Required | Readonly |\n|------|-------|------|----------|----------|\n| `id` | ID | int | no | yes |
| `voucher` | Voucher | int | yes | no |
| `description` | Description | string | yes | no |
| `sale` | Sale | int | no | no |
| `salePayment` | salePayment | int | no | no |
| `payment` | Payment | int | no | no |
| `amount` | Amount | float | yes | no |
| `balanceAfter` | Balance | float | yes | no |
| `unitPrice` | unitPrice | currency | yes | no |
| `useDate` | Production date | datetime | yes | no |
| `fkModel` | fkModel | fkModel | no | yes |
| `fkId` | fkId | int | no | yes |
| `ratio` | Ratio | float | no | no |
| `refunded` | Back | bool | yes | no |
| `type` | Type | select | yes | no |
| `voucherType` | Type | string | yes | no |
| `customerName` | Customer | string | yes | no |
| `beneficiary` | Beneficiary | int | no | no |
| `created` | Creation date | datetime | no | yes |
| `createUser` | Created by | int | no | yes |
| `updated` | Update date | datetime | no | yes |
| `updateUser` | Updated by | int | no | yes |
| `deleted` | Removed | bool | no | yes |

#### Values for `type`

| Value | Name |\n|-------|------|\n| `1` | Recarga |
| `2` | Consumption |
| `3` | Manual |
| `4` | Bulk Top up |


## Voucher

**Table:** `voucher`

### Endpoints\n\n#### Retrieve by ID\n```
GET https://[host]/api/model/voucher/[id]
```

#### Retrieve many\n```
GET https://[host]/api/model/voucher
```

#### Create\n```
POST https://[host]/api/model/voucher
```

#### Update\n```
POST https://[host]/api/model/voucher/[id]
```

#### Delete\n```
DELETE https://[host]/api/model/voucher/[id]
```

### Fields\n\n| Name | Label | Type | Required | Readonly |\n|------|-------|------|----------|----------|\n| `id` | ID | int | no | yes |
| `customer` | Customer | int | yes | no |
| `type` | Type | int | yes | no |
| `unitPrice` | Price per unit | currency | yes | no |
| `balance` | Balance | float | yes | no |
| `startDate` | Valid from | date | no | no |
| `endDate` | Valid until | date | no | no |
| `startTime` | Time from | time | no | no |
| `endTime` | Time until | time | no | no |
| `daysOfWeek` | Days of the week | daysOfWeek | no | no |
| `comments` | Observations | text | no | no |
| `created` | Creation date | datetime | no | yes |
| `createUser` | Created by | int | no | yes |
| `updated` | Update date | datetime | no | yes |
| `updateUser` | Updated by | int | no | yes |
| `deleted` | Removed | bool | no | yes |


## Voucher beneficiary

**Table:** `voucherBeneficiary`

### Endpoints\n\n#### Retrieve by ID\n```
GET https://[host]/api/model/voucherBeneficiary/[id]
```

#### Retrieve many\n```
GET https://[host]/api/model/voucherBeneficiary
```

#### Create\n```
POST https://[host]/api/model/voucherBeneficiary
```

#### Update\n```
POST https://[host]/api/model/voucherBeneficiary/[id]
```

#### Delete\n```
DELETE https://[host]/api/model/voucherBeneficiary/[id]
```

### Fields\n\n| Name | Label | Type | Required | Readonly |\n|------|-------|------|----------|----------|\n| `id` | ID | int | no | yes |
| `customer` | Holder | int | yes | no |
| `beneficiary` | Beneficiary | int | yes | no |
| `created` | Creation date | datetime | no | yes |
| `createUser` | Created by | int | no | yes |
| `updated` | Update date | datetime | no | yes |
| `updateUser` | Updated by | int | no | yes |
| `deleted` | Removed | bool | no | yes |


## Voucher type

**Table:** `voucherType`

### Endpoints\n\n#### Retrieve by ID\n```
GET https://[host]/api/model/voucherType/[id]
```

#### Retrieve many\n```
GET https://[host]/api/model/voucherType
```

#### Create\n```
POST https://[host]/api/model/voucherType
```

#### Update\n```
POST https://[host]/api/model/voucherType/[id]
```

#### Delete\n```
DELETE https://[host]/api/model/voucherType/[id]
```

### Fields\n\n| Name | Label | Type | Required | Readonly |\n|------|-------|------|----------|----------|\n| `id` | ID | int | no | yes |
| `name` | Name | string | yes | no |
| `kind` | Type | select | yes | no |
| `product` | Product | int | yes | no |
| `shortDescription` | Short description | string | no | no |
| `expiration` | Expiration | duration | no | no |
| `startTime` | Time from | time | no | no |
| `endTime` | Time until | time | no | no |
| `daysOfWeek` | Days of the week | daysOfWeek | no | no |
| `customerTag` | Will assign this tag | int | no | no |
| `created` | Creation date | datetime | no | yes |
| `createUser` | Created by | int | no | yes |
| `updated` | Update date | datetime | no | yes |
| `updateUser` | Updated by | int | no | yes |
| `deleted` | Removed | bool | no | yes |

#### Values for `kind`

| Value | Name |\n|-------|------|\n| `1` | Money |
| `3` | Units |


## VoucherPrice

**Table:** `voucherPrice`

### Endpoints\n\n#### Retrieve by ID\n```
GET https://[host]/api/model/voucherPrice/[id]
```

#### Retrieve many\n```
GET https://[host]/api/model/voucherPrice
```

#### Create\n```
POST https://[host]/api/model/voucherPrice
```

#### Update\n```
POST https://[host]/api/model/voucherPrice/[id]
```

#### Delete\n```
DELETE https://[host]/api/model/voucherPrice/[id]
```

### Fields\n\n| Name | Label | Type | Required | Readonly |\n|------|-------|------|----------|----------|\n| `id` | ID | int | no | yes |
| `type` | Type | int | yes | no |
| `tagFilter` | Filter tag | int | no | no |
| `amount` | Top up amount | float | yes | no |
| `price` | Price | currency | yes | no |
| `online` | Available online | bool | yes | no |
| `minAge` | Minimum age | int | no | no |
| `maxAge` | Maximum age | int | no | no |
| `description` | Description | text | no | no |
| `priority` | Priority | int | yes | no |
| `created` | Creation date | datetime | no | yes |
| `createUser` | Created by | int | no | yes |
| `updated` | Update date | datetime | no | yes |
| `updateUser` | Updated by | int | no | yes |
| `deleted` | Removed | bool | no | yes |


