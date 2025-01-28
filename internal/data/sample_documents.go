package data

import "time"

type Document struct {
	Title       string
	Description string
	CreatedAt   time.Time
	Tags        []Tag
}

type Tag struct {
	Name  string
	Color string
}

func GetSampleDocuments() []Document {
	return []Document{
		{
			Title:       "Q4 Financial Report 2023",
			Description: "Comprehensive financial analysis for Q4 2023, including revenue growth, expense analysis, and projections for the upcoming fiscal year.",
			CreatedAt:   time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			Tags: []Tag{
				{Name: "finance", Color: "amber"},
				{Name: "report", Color: "blue"},
			},
		},
		{
			Title:       "Project Roadmap 2024",
			Description: "Strategic planning document outlining key milestones, deliverables, and timeline for major projects in 2024.",
			CreatedAt:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Tags: []Tag{
				{Name: "planning", Color: "green"},
				{Name: "strategy", Color: "purple"},
			},
		},
		{
			Title:       "Meeting Minutes - Product Team",
			Description: "Summary of key discussions and decisions from the monthly product team meeting, including feature prioritization and timeline updates.",
			CreatedAt:   time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
			Tags: []Tag{
				{Name: "meeting", Color: "orange"},
				{Name: "product", Color: "purple"},
			},
		},
		{
			Title:       "Employee Handbook v2.3",
			Description: "Updated company policies and guidelines, including remote work procedures and updated benefits information.",
			CreatedAt:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			Tags: []Tag{
				{Name: "HR", Color: "green"},
				{Name: "policy", Color: "indigo"},
			},
		},
		{
			Title:       "Marketing Campaign Brief",
			Description: "Campaign strategy document for Q2 2024 product launch, including target audience analysis and marketing channels.",
			CreatedAt:   time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC),
			Tags: []Tag{
				{Name: "marketing", Color: "indigo"},
				{Name: "campaign", Color: "purple"},
			},
		},
		{
			Title:       "Technical Documentation - API v3",
			Description: "Comprehensive documentation of the new API endpoints, including authentication methods and response formats.",
			CreatedAt:   time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			Tags: []Tag{
				{Name: "technical", Color: "blue"},
				{Name: "API", Color: "purple"},
			},
		},
		{
			Title:       "Customer Feedback Analysis",
			Description: "Analysis of customer feedback from Q1 2024, including key trends, pain points, and recommended action items.",
			CreatedAt:   time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
			Tags: []Tag{
				{Name: "feedback", Color: "orange"},
				{Name: "analysis", Color: "green"},
			},
		},
		{
			Title:       "Security Protocol Guidelines",
			Description: "Updated security protocols and best practices for internal systems and data protection procedures.",
			CreatedAt:   time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
			Tags: []Tag{
				{Name: "security", Color: "red"},
				{Name: "protocol", Color: "green"},
			},
		},
		{
			Title:       "Budget Forecast 2024",
			Description: "Detailed budget projections for Q3 and Q4 2024, including department allocations and investment plans.",
			CreatedAt:   time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC),
			Tags: []Tag{
				{Name: "finance", Color: "teal"},
				{Name: "forecast", Color: "blue"},
			},
		},
		{
			Title:       "Product Launch Plan",
			Description: "Detailed execution plan for the upcoming product launch, including marketing strategy and release timeline.",
			CreatedAt:   time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
			Tags: []Tag{
				{Name: "product", Color: "red"},
				{Name: "launch", Color: "pink"},
				{Name: "planning", Color: "green"},
			},
		},
	}
}
