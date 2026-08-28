package utils

const ReferralFormTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Referral Request</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Montserrat:ital,wght@0,100..900;1,100..900&display=swap" rel="stylesheet">
    <style>
        @page {
        @bottom-left {
            font-size: 9px;
            color: #6e6e6e;
            font-family: 'Montserrat';
            white-space: pre-line;
            margin-top: -29px;
            margin-left: 4px;
        }
        size: A4;
        margin: 0;
        }
        html{
            margin: 0;
        }
        body {
            font-family: Montserrat;
            font-size: 12px;
            margin: 0;
        }
        .page{
            display: flex;
            justify-content: center;
            padding: 0 30px;
        }
        .content{
            max-width: 900px;
            width: calc(100% - 30px);
        }
        hr{
            background-color: #981098;
            height: 1px;
            border: 0
        }
        h3{
            margin-bottom: 10px;
            font-weight: 600;
            font-size: 13px;
        }
        .header {
            display: flex;
            align-items: flex-end;
            justify-content: start; /* This pushes the logo and details to opposite ends */
            margin-bottom: 20px;
            margin-top: 20px;
            padding: 0 20px;
        }
        .header-logo {
            flex-grow: 0;
            flex-shrink: 0;
        }
        .header-logo img {
            height: 50px;
            width: auto;
        }
        .header-logo, .header-details {
            display: flex;
            align-items: start;
        }
        .header-details {
            padding-left: 20px;
            justify-content: flex-start; /* Aligns details to the end (right) */
            text-align: right;
        }
        .header-title {
            color: #800080;
            font-weight: bold;
            font-size: 25px;
            margin: 0 auto;
            text-align: center;
            flex: 2;
        }
        .report-title {
            font-size: 18px;
            font-weight: 600;
            margin: 20px 0;
        }
        .detail-section {
            margin-bottom: 30px;
            border-radius: 8px;
            display: flex;
            flex-wrap: wrap;
            padding: 15px;
            align-items: start;
            border: 2px solid #ebebeb;
        }
        .detail-section > .hd-section.w-100:first-child{
            margin: 0;
        }
        .detail-section > .hd-section.w-50:first-child{
            margin: 0;
        }
        .detail-section > .hd-section.w-50:nth-child(2){
            margin: 0;
        }
        .hd-section {
            margin: 16px 0 0;
            display: flex;
            justify-content: start;
        }
        .w-50{
            width: 50%;
        }
        .w-100{
            width: 100%;
        }
        .hd-title {
            color: #838282;
            min-width: 120px;
            font-weight: 500;
            font-size: 11px;
            padding-bottom: 2px;
        }
        .hd-value {
            font-size: 11px;
            font-weight: 500;
        }
        table {
            margin-top: 15px;
            width: 100%;
            border-collapse: collapse;
            font-size: 11px;
        }
        th{
            color: #838282;
            font-weight: 500;
        }
        td{
            font-weight: 500;
        }
        tr > th:first-child, tr > td:first-child{
            padding-left: 0;
        }
        th, td {
            padding: 8px;
            text-align: left;
            border: none;
        }
        .footer {
            font-size: 9px;
            max-width: 950px;
            text-align: center;
            padding: 20px;
            color: #4d184d;
            background-color: #4d184d33;
            margin-top: 5px;
            font-weight: 700;
            display: flex;
            justify-content: space-around;
            align-items: center;
        }
        .footer > .flex-row{
            display: flex;
            flex-direction: row;
        }

        .flex-row > .ft-title{
            font-weight: 500;
            margin-right: 5px;
        }
        .flex-row > .ft-divider{
            width: 5px;
        }
    </style>
</head>
<body>
    <div class="page">
        <div class="content">
            <div class="header">
                <div class="header-logo">
                    <img src="https://storage.googleapis.com/mycarehub/Empower%20Logo.png" alt="Logo" height="70px">
                </div>
                <div class="header-details">
                    <h1 class="header-title">{{ .Facility }}</h1>
                </div>
            </div>
            <hr />
            <div class="report-title">Referral Request</div>

            {{if .Patient}}
            <h3>Patient information</h3>
            <div class="detail-section">
                {{if .Patient.Name}}
                <div class="hd-section w-50">
                    <div class="hd-title">Name</div>
                    <div class="hd-value">{{.Patient.Name}}</div>
                </div>
                {{end}}
                {{if .Patient.DateOfBirth}}
                <div class="hd-section w-50">
                    <div class="hd-title">Date of birth</div>
                    <div class="hd-value">{{.Patient.DateOfBirth}} ({{.Patient.Age}})</div>
                </div>
                {{end}}
                {{if .Patient.EmpowerID}}
                <div class="hd-section w-50">
                    <div class="hd-title">Empower ID</div>
                    <div class="hd-value">{{.Patient.EmpowerID}}</div>
                </div>
                {{end}}
                {{if .Patient.Sex}}
                <div class="hd-section w-50">
                    <div class="hd-title">Sex</div>
                    <div class="hd-value">{{.Patient.Sex}}</div>
                </div>
                {{end}}
                {{if .Patient.NationalID}}
                <div class="hd-section w-50">
                    <div class="hd-title">National ID</div>
                    <div class="hd-value">{{.Patient.NationalID}}</div>
                </div>
                {{end}}
                {{if .Patient.PhoneNumber}}
                <div class="hd-section w-50">
                    <div class="hd-title">Phone number</div>
                    <div class="hd-value">{{.Patient.PhoneNumber}}</div>
                </div>
                {{end}}
            </div>
            {{end}}

            {{if .ReferralReason}}
            <h3>Referral details</h3>
            <div class="detail-section">
                {{if .ReferralReason.Reason}}
                <div class="hd-section w-100">
                    <div class="hd-title">Reason</div>
                    <div class="hd-value">{{.ReferralReason.Reason}}</div>
                </div>
                {{end}}
                {{if .ReferralReason.Test}}
                <div class="hd-section w-100">
                    <div class="hd-title">Test</div>
                    <div class="hd-value">{{.ReferralReason.Test}}</div>
                </div>
                {{end}}
                {{if .ReferralReason.Date}}
                <div class="hd-section w-100">
                    <div class="hd-title">Referral date</div>
                    <div class="hd-value">{{.ReferralReason.Date}}</div>
                </div>
                {{end}}
                {{if .ReferralReason.Note}}
                <div class="hd-section w-100">
                    <div class="hd-title">Notes</div>
                    <div class="hd-value">{{.ReferralReason.Note}}</div>
                </div>
                {{end}}
            </div>
            {{end}}

            {{if or .ReceivingFacility.FacilityName .ReceivingFacility.FacilityContact .ReceivingFacility.FacilityCounty}}
            <h3>Receiving facility</h3>
            <div class="detail-section">
                {{if .ReceivingFacility.FacilityName}}
                <div class="hd-section w-50">
                    <div class="hd-title">Referred to</div>
                    <div class="hd-value">{{.ReceivingFacility.FacilityName}}</div>
                </div>
                {{end}}
                {{if .ReceivingFacility.FacilityContact}}
                <div class="hd-section w-50">
                    <div class="hd-title">Hospital contact</div>
                    <div class="hd-value">{{.ReceivingFacility.FacilityContact}}</div>
                </div>
                {{end}}
                {{if .ReceivingFacility.FacilityCounty}}
                <div class="hd-section w-100">
                    <div class="hd-title">Location</div>
                    <div class="hd-value">{{.ReceivingFacility.FacilityCounty}}</div>
                </div>
                {{end}}
            </div>
            {{end}}

            {{if or .NextOfKin.Name .NextOfKin.PhoneNumber .NextOfKin.Relationship}}
            <h3>Next of kin</h3>
            <div class="detail-section">
                {{if .NextOfKin.Name}}
                <div class="hd-section w-50">
                    <div class="hd-title">Name</div>
                    <div class="hd-value">{{.NextOfKin.Name}}</div>
                </div>
                {{end}}
                {{if .NextOfKin.PhoneNumber}}
                <div class="hd-section w-50">
                    <div class="hd-title">Phone number</div>
                    <div class="hd-value">{{.NextOfKin.PhoneNumber}}</div>
                </div>
                {{end}}
                {{if .NextOfKin.Relationship}}
                <div class="hd-section w-50">
                    <div class="hd-title">Relationship</div>
                    <div class="hd-value">{{.NextOfKin.Relationship}}</div>
                </div>
                {{end}}
            </div>
            {{end}}

            {{if .MedicalHistory}}
            <h3>Medical history</h3>
            <div class="detail-section">
                {{if .MedicalHistory.Procedure}}
                <div class="hd-section w-100">
                    <div class="hd-title">Procedure</div>
                    <div class="hd-value">{{.MedicalHistory.Procedure}}</div>
                </div>
                {{end}}
                {{if .MedicalHistory.Medication}}
                <div class="hd-section w-100">
                    <div class="hd-title">Medication</div>
                    <div class="hd-value">{{.MedicalHistory.Medication}}</div>
                </div>
                {{end}}
                {{if .MedicalHistory.Tests}}
                <div class="w-100">
                    <table>
                        <thead>
                            <tr>
                                <th>Test Done</th>
                                <th>Results</th>
                                <th>Date</th>
                            </tr>
                        </thead>
                        <tbody>
                            {{range .MedicalHistory.Tests}}
                            <tr>
                                {{if .Name}}<td>{{.Name}}</td>{{end}}
                                {{if .Results}}<td>{{.Results}}</td>{{else}}<td></td>{{end}}
                                {{if .Date}}<td>{{.Date}}</td>{{else}}<td></td>{{end}}
                            </tr>
                            {{end}}
                        </tbody>
                    </table>
                </div>
                {{end}}
            </div>
            {{end}}

            <h3>Referred by</h3>
            <div class="detail-section">
                <div class="hd-section w-50">
                    <div class="hd-title">Referring Officer</div>
                    <div class="hd-value"></div>
                </div>
                <div class="hd-section w-50">
                    <div class="hd-title">Designation</div>
                    <div class="hd-value"></div>
                </div>
                <div class="hd-section w-50">
                    <div class="hd-title">Phone</div>
                    <div class="hd-value"></div>
                </div>
                <div class="hd-section w-50">
                    <div class="hd-title">Signature</div>
                    <div class="hd-value"></div>
                </div>
            </div>
        </div>
    </div>

    <div class="footer">
        <div class="flex-row">
            <div class="ft-title">
                Phone:
            </div>
            <div>
                +254720999888
            </div>
        </div>
        <div class="ft-divider"> | </div>
        <div class="flex-row">
            <div class="ft-title">Email:</div>
            <div>info@nshospital.com</div>
        </div>
        <div class="ft-divider"> | </div>
        <div class="flex-row">
            <div class="ft-title">
                Postal Address:
            </div>
            <div>
                PO Box 1234-00100 Nairobi
            </div>
        </div>
    </div>
</body>
</html>
`

const ReferralEmailTemplate = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Email Template</title>
    <style>
      body {
        font-family: Montserrat, Arial, sans-serif;
        margin: 0;
        padding: 0;
        background-color: #f4f4f4;
      }
      .header {
        background-color: #4d184d;
        color: white;
        padding: 20px 40px;
        text-align: left;
        box-sizing: border-box;
      }
      .header h1 {
        margin: 0;
        font-size: 28px;
        font-weight: 700;
        line-height: 45.27px;
        word-wrap: break-word;
      }
      .header p {
        margin: 5px 0;
        font-size: 16px;
        line-height: 26px;
        word-wrap: break-word;
        font-weight: 300;
      }
      .header p:first-of-type {
        font-weight: 500;
      }
      .header p span {
        font-weight: 600;
      }
      .header p:last-of-type {
        font-weight: 400;
      }
      .container {
        background-color: #ffffff;
        padding: 20px 40px;
        box-sizing: border-box;
        width: 100%;
        border: 0.5px #cecece solid;
      }
      .content p {
        margin: 10px 0;
      }
      .details {
        margin: 20px 0;
      }
      .details p {
        margin: 5px 0;
        color: #666666;
        font-size: 16px;
        line-height: 25px;
      }
      .details p span {
        font-weight: 600;
      }
      .footer {
        text-align: left;
        color: #666666;
        font-size: 12px;
        margin-top: 20px;
      }
      .footer p {
        margin: 10px 0;
        font-size: 16px;
        font-weight: 500;
        line-height: 25px;
      }
      .footer-note {
        color: #666666;
        font-size: 12px;
        font-weight: 400;
      }

      .from-empower {
        font-size: 16px;
        font-weight: 550;
        margin-top: 10px;
        margin-bottom: 20px;
        line-height: 25px;
      }

      .details-title {
        color: #666666;
        font-size: 16px;
        font-weight: 600;
        line-height: 25px;
      }
    </style>
  </head>
  <body>
    <div class="header">
      <img
        src="https://storage.googleapis.com/mycarehub/empower-logo.png"
        alt="Empower Logo"
        height="100px"
      />
      <p>
        Hello, you have received a new referral request from
        <span>Empower Makueni</span>. Please find the referral report attached.
        Here are the details:
      </p>
    </div>
    <div class="container">
      <div class="content">
        <div class="details">
          <p class="details-title">Details</p>
          <p>Patient Name: <span>{{.PatientName}}</span></p>
          <p>Phone Number: <span>{{.PatientPhoneNumber}}</span></p>
        </div>
        <p>Warm regards,</p>
        <p class="from-empower">Empower</p>
        <p class="footer-note">
          NB: This is an automated email. Do not reply to it.
        </p>
      </div>
    </div>
  </body>
</html>
`
