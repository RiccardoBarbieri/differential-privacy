package model

import (
	"encoding/csv"
	"fmt"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/register"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const AdmissionColNumber = 15

func init() {
	register.Function2x1[string, func(admission Admission), error](CreateAdmissionFn)
	register.Emitter1[Admission]()

}

// Admission represents the event of an admission to a hospital of a patient
type Admission struct {
	Name              string
	Age               int
	Gender            string
	BloodType         string
	MedicalCondition  string
	DateOfAdmission   time.Time
	Doctor            string
	Hospital          string
	InsuranceProvider string
	BillingAmount     float64
	RoomNumber        int
	AdmissionType     string
	DischargeDate     time.Time
	Medication        string
	TestResults       string
}

func CreateAdmissionFn(line string, emit func(admission Admission)) error {
	// notHeader is true if the line contains a number (no numbers in header line)
	notHeader, err := regexp.MatchString("[0-9]", line)
	if err != nil {
		return err
	}
	if !notHeader {
		return nil
	}

	admission := Admission{}
	reader := csv.NewReader(strings.NewReader(line))
	cols, err := reader.Read()
	if err != nil {
		return err
	}
	if len(cols) != AdmissionColNumber {
		return fmt.Errorf("line containse %d columns, Admissions expects %d columns - %s", len(cols), AdmissionColNumber, line)
	}
	admission.Name = cols[0]
	admission.Age, err = strconv.Atoi(cols[1])

	admission.Gender = cols[2]
	admission.BloodType = cols[3]
	admission.MedicalCondition = cols[4]
	admission.DateOfAdmission, err = time.Parse(time.DateOnly, cols[5])
	if err != nil {
		return err
	}
	admission.Doctor = cols[6]
	admission.Hospital = cols[7]
	admission.InsuranceProvider = cols[8]
	admission.BillingAmount, err = strconv.ParseFloat(cols[9], 64)
	if err != nil {
		return err
	}
	admission.RoomNumber, err = strconv.Atoi(cols[10])
	if err != nil {
		return err
	}
	admission.AdmissionType = cols[11]
	admission.DischargeDate, err = time.Parse(time.DateOnly, cols[12])
	if err != nil {
		return err
	}
	admission.Medication = cols[13]
	admission.TestResults = cols[14]

	emit(admission)
	return nil
}

func (a Admission) String() string {
	return fmt.Sprintf("Name = %s, Admission = %s, Discharge %s", a.Name, a.DateOfAdmission.Format(time.RFC3339), a.DischargeDate.Format(time.RFC3339))
}
