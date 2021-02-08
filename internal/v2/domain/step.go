package domain

type Step int

const (
	Step1 Step = iota + 1
	Step2
	Step3
	Step4
	Step5
	Step6
	Step7
	Step8
	Step9
	Step10
	//StepCount marks the the length of possible steps (StepCount-1 == last possible step)
	StepCount
)
